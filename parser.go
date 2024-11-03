package jsontools

import (
	"errors"
	"fmt"
)

type Kind byte

const (
	KindOther Kind = iota
	KindObjectKey
	KindObjectValue
	KindArrayValue
)

func (s Kind) String() string {
	switch s {
	case KindObjectKey:
		return "key"
	case KindObjectValue:
		return "value"
	case KindArrayValue:
		return "array"
	default:
		return "unknown"
	}
}

type jsonParser struct {
	jsonTokenizer
	stack       []Kind
	userHandler jsonParserHandler
	err         error
}

func NewJsonParser(data []byte, handler jsonParserHandler) *jsonParser {
	return &jsonParser{
		jsonTokenizer: jsonTokenizer{
			data:    data,
			current: Init,
		},
		stack:       make([]Kind, 0, 32),
		userHandler: handler,
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

type HandlerContext struct {
	Token     TokenType
	Kind      Kind
	Value     []byte
	StackSize int
}

type jsonParserHandler func(ctx HandlerContext) error

func (t *jsonParser) push(kind Kind) {
	t.stack = append(t.stack, kind)
}

func (t *jsonParser) pop() Kind {
	kind := t.stack[len(t.stack)-1]
	t.stack = t.stack[:len(t.stack)-1]
	return kind
}

func (t *jsonParser) peek() Kind {
	return t.stack[len(t.stack)-1]
}

func (t *jsonParser) isEmpty() bool {
	return len(t.stack) == 0
}

func (t *jsonParser) handler(token TokenType, kind Kind, value []byte) {
	t.err = t.userHandler(HandlerContext{
		Token:     token,
		Kind:      kind,
		Value:     value,
		StackSize: len(t.stack),
	})
}

func (t *jsonParser) Parse() error {
	flags := flagBeginObject | flagBeginArray
	for {
		if t.err != nil {
			return t.err
		}

		token, value, err := t.Next()
		if err != nil {
			return err
		}

		switch token {
		case BeginObject:
			if !flags.has(flagBeginObject) {
				return errors.New("invalid '{'")
			}
			t.push(KindObjectValue)
			t.handler(token, KindOther, value)
			// current is '{', expect '}' or key
			flags = flagObjectKey | flagBeginObject | flagEndObject

		case String:
			if flags.has(flagObjectKey) {
				// current is key, expect ':'
				t.push(KindObjectKey)
				t.handler(token, KindObjectKey, value)
				flags = flagColon
				continue
			}
			if flags.has(flagObjectValue) {
				// current is object value, expect ',' or '}'
				t.handler(token, KindObjectValue, value)
				t.pop()
				flags = flagComma | flagEndObject
				continue
			}
			if flags.has(flagArrayValue) {
				// current is array value, expect ',' or ']'
				t.handler(token, KindArrayValue, value)
				flags = flagComma | flagEndArray
				continue
			}
			return fmt.Errorf("invalid string '%s'", string(value))

		case Number:
			if flags.has(flagObjectValue) {
				// current is object value, expect ',' or '}'
				t.handler(token, KindObjectValue, value)
				t.pop()
				flags = flagComma | flagEndObject
				continue
			}
			if flags.has(flagArrayValue) {
				// current is array value, expect ',' or ']'
				t.handler(token, KindArrayValue, value)
				flags = flagComma | flagEndArray
				continue
			}
			return fmt.Errorf("invalid int '%s'", string(value))

		case Float:
			if flags.has(flagObjectValue) {
				// current is object value, expect ',' or '}'
				t.handler(token, KindObjectValue, value)
				t.pop()
				flags = flagComma | flagEndObject
				continue
			}
			if flags.has(flagArrayValue) {
				// current is array value, expect ',' or ']'
				t.handler(token, KindArrayValue, value)
				flags = flagComma | flagEndArray
				continue
			}
			return fmt.Errorf("invalid float '%s'", string(value))

		case True:
			if flags.has(flagObjectValue) {
				// current is object value, expect ',' or '}'
				t.handler(token, KindObjectValue, value)
				t.pop()
				flags = flagComma | flagEndObject
				continue
			}
			if flags.has(flagArrayValue) {
				// current is array value, expect ',' or ']'
				t.handler(token, KindArrayValue, value)
				flags = flagComma | flagEndArray
				continue
			}
			return fmt.Errorf("invalid bool '%s'", string(value))

		case False:
			if flags.has(flagObjectValue) {
				// current is object value, expect ',' or '}'
				t.handler(token, KindObjectValue, value)
				t.pop()
				flags = flagComma | flagEndObject
				continue
			}
			if flags.has(flagArrayValue) {
				// current is array value, expect ',' or ']'
				t.handler(token, KindArrayValue, value)
				flags = flagComma | flagEndArray
				continue
			}
			return fmt.Errorf("invalid bool '%s'", string(value))

		case Null:
			if flags.has(flagObjectValue) {
				// current is object value, expect ',' or '}'
				t.handler(token, KindObjectValue, value)
				t.pop()
				flags = flagComma | flagEndObject
				continue
			}
			if flags.has(flagArrayValue) {
				// current is array value, expect ',' or ']'
				t.handler(token, KindArrayValue, value)
				flags = flagComma | flagEndArray
				continue
			}
			return fmt.Errorf("invalid null '%s'", string(value))

		case SepComma:
			if !flags.has(flagComma) {
				return errors.New("invalid ','")
			}
			if flags.has(flagEndObject) {
				// current is ',' in object, expect key
				t.handler(token, KindOther, value)
				flags = flagObjectKey
				continue
			}
			if flags.has(flagEndArray) {
				// current is ',' in array, expect value or object
				t.handler(token, KindOther, value)
				flags = flagArrayValue | flagBeginArray | flagBeginObject
				continue
			}

		case SepColon:
			if !flags.has(flagColon) {
				return errors.New("invalid ':'")
			}
			// current is ':', expect value, object or array
			t.handler(token, KindOther, value)
			flags = flagObjectValue | flagBeginObject | flagBeginArray

		case BeginArray:
			if flags.has(flagBeginArray) {
				// current is '[', expect ']' or value or object
				t.push(KindArrayValue)
				t.handler(token, KindOther, value)
				flags = flagArrayValue | flagBeginArray | flagEndArray | flagBeginObject
			}

		case EndArray:
			if !flags.has(flagEndArray) {
				return errors.New("invalid ']'")
			}
			t.pop()
			if t.isEmpty() {
				// current is ']', but if stack is empty, which means array is top level.
				// expect another object/array.
				// eg: "[1, 2, 3]"
				t.push(KindArrayValue)
				t.handler(token, KindOther, value)
				flags = flagBeginObject | flagBeginArray
				continue
			}

			kind := t.peek()
			if kind == KindObjectKey {
				// current is ']', but if there's key in stack, which means array is inside object.
				// expect ',' or '}'
				// eg: {"arr":[1, 2, 3], "nextKey": "value"}
				t.handler(token, KindOther, value)
				t.pop()
				flags = flagComma | flagEndObject
				continue
			}

			if kind == KindArrayValue {
				// current is ']', but if there's array in stack, which means array is inside array.
				// expect ',' or ']'
				// eg: [1, 2, 3, [4, 5, 6], 7]
				t.handler(token, KindOther, value)
				flags = flagComma | flagEndArray
				continue
			}

			return errors.New("missing ']'")

		case EndObject:
			if !flags.has(flagEndObject) {
				return errors.New("invalid '}'")
			}
			t.pop()
			if t.isEmpty() {
				// current is '}', but if stack is empty, which means object is top level.
				// expect another object/array.
				// eg: "[1, 2, 3]"
				t.push(KindObjectValue)
				t.handler(token, KindOther, value)
				flags = flagBeginObject | flagBeginArray
				continue
			}

			kind := t.peek()
			if kind == KindObjectKey {
				// current is '}', but if there's key in stack, which means object is inside object.
				// expect ',' or '}'
				// eg: {"obj":{"key1":1, "key2":2}, "nextKey": "value"}
				t.handler(token, KindOther, value)
				t.pop()
				flags = flagComma | flagEndObject
				continue
			}

			if kind == KindArrayValue {
				// current is '}', but if there's array in stack, which means object is inside array.
				// expect ',' or ']'
				// eg: [1, 2, 3, {"key1":1, "key2":2}, 7]
				t.handler(token, KindOther, value)
				flags = flagComma | flagEndArray
				continue
			}

			return errors.New("invalid '}'")

		case EndJson:
			if t.isEmpty() {
				return errors.New("invalid EOF1")
			}

			t.pop()
			if t.isEmpty() {
				return nil
			}
			return errors.New("invalid EOF2")
		}
	}
}
