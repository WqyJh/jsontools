package jsontools

import (
	"errors"
	"fmt"
)

type StackType byte

const (
	StackTypeValue StackType = iota
	StackTypeObjectKey
	StackTypeArray
)

func (s StackType) String() string {
	switch s {
	case StackTypeValue:
		return "value"
	case StackTypeObjectKey:
		return "key"
	case StackTypeArray:
		return "array"
	default:
		return "unknown"
	}
}

type jsonParser struct {
	jsonTokenizer
	stack   []StackType
	handler jsonParserHandler
}

func NewJsonParser(data []byte, handler jsonParserHandler) *jsonParser {
	return &jsonParser{
		jsonTokenizer: *NewJsonTokenizer(data),
		stack:         make([]StackType, 0, 128),
		handler:       handler,
	}
}

type parseFlags int

const (
	flagObjectKey   parseFlags = 0x0001
	flagColon       parseFlags = 0x0002
	flagObjectValue parseFlags = 0x0004
	flagComma       parseFlags = 0x0008
	flagBeginObject parseFlags = 0x0010
	flagEndObject   parseFlags = 0x0020
	flagBeginArray  parseFlags = 0x0040
	flagEndArray    parseFlags = 0x0080
	flagArrayValue  parseFlags = 0x0100
)

func (f parseFlags) has(flag parseFlags) bool {
	return f&flag != 0
}

type jsonParserHandler func(token TokenType, value []byte)

func (t *jsonParser) push(stackType StackType) {
	t.stack = append(t.stack, stackType)
}

func (t *jsonParser) pop() StackType {
	stackType := t.stack[len(t.stack)-1]
	t.stack = t.stack[:len(t.stack)-1]
	return stackType
}

func (t *jsonParser) peek() StackType {
	return t.stack[len(t.stack)-1]
}

func (t *jsonParser) isEmpty() bool {
	return len(t.stack) == 0
}

func (t *jsonParser) Parse() error {
	flags := flagBeginObject | flagBeginArray
	for {
		token, value, err := t.Next()
		if err != nil {
			return err
		}

		if token == EndJson {
			t.pop()
			if t.isEmpty() {
				return nil
			}
			return errors.New("invalid EOF")
		}

		t.handler(token, value)

		switch token {
		case BeginObject:
			if !flags.has(flagBeginObject) {
				return errors.New("invalid '{'")
			}
			t.push(StackTypeValue)
			// current is '{', expect '}' or key
			flags = flagObjectKey | flagBeginObject | flagEndObject

		case String:
			if flags.has(flagObjectKey) {
				// current is key, expect ':'
				t.push(StackTypeObjectKey)
				flags = flagColon
				continue
			}
			if flags.has(flagObjectValue) {
				// current is object value, expect ',' or '}'
				t.pop()
				flags = flagComma | flagEndObject
				continue
			}
			if flags.has(flagArrayValue) {
				// current is array value, expect ',' or ']'
				flags = flagComma | flagEndArray
				continue
			}
			return fmt.Errorf("invalid string '%s'", string(t.value))

		case Number:
			if flags.has(flagObjectValue) {
				// current is object value, expect ',' or '}'
				t.pop()
				flags = flagComma | flagEndObject
				continue
			}
			if flags.has(flagArrayValue) {
				// current is array value, expect ',' or ']'
				flags = flagComma | flagEndArray
				continue
			}
			return fmt.Errorf("invalid int '%s'", string(t.value))

		case Float:
			if flags.has(flagObjectValue) {
				// current is object value, expect ',' or '}'
				t.pop()
				flags = flagComma | flagEndObject
				continue
			}
			if flags.has(flagArrayValue) {
				// current is array value, expect ',' or ']'
				flags = flagComma | flagEndArray
				continue
			}
			return fmt.Errorf("invalid float '%s'", string(t.value))

		case True:
			if flags.has(flagObjectValue) {
				// current is object value, expect ',' or '}'
				t.pop()
				flags = flagComma | flagEndObject
				continue
			}
			if flags.has(flagArrayValue) {
				// current is array value, expect ',' or ']'
				flags = flagComma | flagEndArray
				continue
			}
			return fmt.Errorf("invalid bool true '%s'", string(t.value))

		case False:
			if flags.has(flagObjectValue) {
				// current is object value, expect ',' or '}'
				t.pop()
				flags = flagComma | flagEndObject
				continue
			}
			if flags.has(flagArrayValue) {
				// current is array value, expect ',' or ']'
				flags = flagComma | flagEndArray
				continue
			}
			return fmt.Errorf("invalid bool false '%s'", string(t.value))

		case Null:
			if flags.has(flagObjectValue) {
				// current is object value, expect ',' or '}'
				t.pop()
				flags = flagComma | flagEndObject
				continue
			}
			if flags.has(flagArrayValue) {
				// current is array value, expect ',' or ']'
				flags = flagComma | flagEndArray
				continue
			}
			return fmt.Errorf("invalid null '%s'", string(t.value))

		case SepComma:
			if !flags.has(flagComma) {
				return errors.New("missing ','")
			}
			if flags.has(flagEndObject) {
				// current is ',' in object, expect key
				flags = flagObjectKey
				continue
			}
			if flags.has(flagEndArray) {
				// current is ',' in array, expect value or object
				flags = flagArrayValue | flagBeginArray | flagBeginObject
				continue
			}

		case SepColon:
			if !flags.has(flagColon) {
				return errors.New("missing ':'")
			}
			// current is ':', expect value, object or array
			flags = flagObjectValue | flagBeginObject | flagBeginArray

		case BeginArray:
			if flags.has(flagBeginArray) {
				t.push(StackTypeArray)
				// current is '[', expect ']' or value or object
				flags = flagArrayValue | flagBeginArray | flagEndArray | flagBeginObject
			}

		case EndArray:
			if !flags.has(flagEndArray) {
				return errors.New("invalid ']'")
			}
			t.pop()
			if t.isEmpty() {
				t.push(StackTypeArray)
				continue
			}

			stackType := t.peek()
			if stackType == StackTypeObjectKey {
				// current is ']', but if there's key in stack, which means array is inside object.
				// expect ',' or '}'
				// eg: {"arr":[1, 2, 3], "nextKey": "value"}
				t.pop()
				flags = flagComma | flagEndObject
				continue
			}

			if stackType == StackTypeArray {
				// current is ']', but if there's array in stack, which means array is inside array.
				// expect ',' or ']'
				// eg: [1, 2, 3, [4, 5, 6], 7]
				flags = flagComma | flagEndArray
				continue
			}

			return fmt.Errorf("missing ']'")

		case EndObject:
			if !flags.has(flagEndObject) {
				return errors.New("invalid '}'")
			}
			t.pop()
			if t.isEmpty() {
				t.push(StackTypeValue)
				continue
			}

			stackType := t.peek()
			if stackType == StackTypeObjectKey {
				// current is '}', but if there's key in stack, which means object is inside object.
				// expect ',' or '}'
				// eg: {"obj":{"key1":1, "key2":2}, "nextKey": "value"}
				t.pop()
				flags = flagComma | flagEndObject
				continue
			}

			if stackType == StackTypeArray {
				// current is '}', but if there's array in stack, which means object is inside array.
				// expect ',' or ']'
				// eg: [1, 2, 3, {"key1":1, "key2":2}, 7]
				flags = flagComma | flagEndArray
				continue
			}

			return fmt.Errorf("invalid '}'")
		}
	}
}
