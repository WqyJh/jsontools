package jsontools

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

type TokenType byte

const (
	Init        TokenType = iota
	BeginObject           // {
	EndObject             // }
	BeginArray            // [
	EndArray              // ]
	Null                  // null
	Number                // number
	Float                 // float
	String                // "string"
	True                  // true
	False                 // false
	SepColon              // :
	SepComma              // ,
	EndJson               // EOF
)

func (t TokenType) String() string {
	switch t {
	case Init:
		return "Init"
	case BeginObject:
		return "BeginObject"
	case EndObject:
		return "EndObject"
	case BeginArray:
		return "BeginArray"
	case EndArray:
		return "EndArray"
	case Null:
		return "Null"
	case Number:
		return "Number"
	case Float:
		return "Float"
	case String:
		return "String"
	case True:
		return "True"
	case False:
		return "False"
	case SepColon:
		return "SepColon"
	case SepComma:
		return "SepComma"
	case EndJson:
		return "EndJson"
	default:
		return "Unknown"
	}
}

type jsonTokenizer struct {
	data []byte

	off     int
	current TokenType
	start   int    // token start
	value   []byte // token value
}

func NewJsonTokenizer(data []byte) *jsonTokenizer {
	return &jsonTokenizer{
		data:    data,
		current: Init,
	}
}

func isDigit(b rune, includeSign bool) bool {
	if b >= 48 && b <= 57 {
		return true
	}

	return includeSign && b == '-'
}

func (t *jsonTokenizer) nextStatus(b rune, size int) TokenType {
	switch b {
	case '{':
		t.value = t.data[t.off : t.off+size]
		return BeginObject
	case '}':
		t.value = t.data[t.off : t.off+size]
		return EndObject
	case '[':
		t.value = t.data[t.off : t.off+size]
		return BeginArray
	case ']':
		t.value = t.data[t.off : t.off+size]
		return EndArray
	case 'n':
		t.start = t.off
		return Null
	case 't':
		t.start = t.off
		return True
	case 'f':
		t.start = t.off
		return False
	case ':':
		t.value = t.data[t.off : t.off+size]
		return SepColon
	case ',':
		t.value = t.data[t.off : t.off+size]
		return SepComma
	case '"':
		t.start = t.off
		return String
	}
	if isDigit(b, true) {
		t.start = t.off
		return Number
	}
	return Init
}

func (t *jsonTokenizer) pendingNextStatus() TokenType {
	return Init
}

func (t *jsonTokenizer) Next() (TokenType, []byte, error) {
	for t.off < len(t.data) {
		b, size := utf8.DecodeRune(t.data[t.off:])

		switch t.current {
		case Init:
			t.current = t.nextStatus(b, size)
			t.off += size

		case BeginObject:
			value := t.value
			t.current = t.nextStatus(b, size)
			t.off += size
			return BeginObject, value, nil

		case EndObject:
			value := t.value
			t.current = t.nextStatus(b, size)
			t.off += size
			return EndObject, value, nil

		case String:
			slashCount := 0
			for j := t.off; j < len(t.data); {
				b, size := utf8.DecodeRune(t.data[j:])
				switch b {
				case '\\':
					slashCount++
				case '"':
					if slashCount%2 == 0 {
						value := t.data[t.start : j+size]
						t.current = t.pendingNextStatus()
						t.off = j + size
						return String, value, nil
					}
				default:
					slashCount = 0
				}
				j += size
			}

			return Init, nil, fmt.Errorf("invalid string '%s'", string(t.data[t.start:]))

		case Number:
			if b == '.' {
				t.current = Float
				t.off += size
			} else if isDigit(b, false) {
				t.off += size
			} else {
				value := t.data[t.start:t.off]
				t.current = t.nextStatus(b, size)
				t.off += size
				return Number, value, nil
			}

		case Float:
			if b == '.' {
				return Init, nil, errors.New("invalid float")
			}
			if isDigit(b, false) {
				t.off += size
			} else {
				value := t.data[t.start:t.off]
				t.current = t.nextStatus(b, size)
				t.off += size
				return Float, value, nil
			}

		case SepColon:
			value := t.value
			t.current = t.nextStatus(b, size)
			t.off += size
			return SepColon, value, nil

		case SepComma:
			value := t.value
			t.current = t.nextStatus(b, size)
			t.off += size
			return SepComma, value, nil

		case BeginArray:
			value := t.value
			t.current = t.nextStatus(b, size)
			t.off += size
			return BeginArray, value, nil

		case EndArray:
			value := t.value
			t.current = t.nextStatus(b, size)
			t.off += size
			return EndArray, value, nil

		case True:
			if t.off+3 <= len(t.data) && t.data[t.off] == 'r' && t.data[t.off+1] == 'u' && t.data[t.off+2] == 'e' {
				value := t.data[t.start : t.off+3]
				t.off += 3
				t.current = t.pendingNextStatus()
				return True, value, nil
			}
			return Init, nil, fmt.Errorf("invalid bool true '%s'", string(t.data[t.start:]))

		case False:
			if t.off+4 <= len(t.data) && t.data[t.off] == 'a' && t.data[t.off+1] == 'l' && t.data[t.off+2] == 's' && t.data[t.off+3] == 'e' {
				value := t.data[t.start : t.off+4]
				t.off += 4
				t.current = t.pendingNextStatus()
				return False, value, nil
			}
			return Init, nil, fmt.Errorf("invalid bool false '%s'", string(t.data[t.start:]))

		case Null:
			if t.off+3 <= len(t.data) && t.data[t.off] == 'u' && t.data[t.off+1] == 'l' && t.data[t.off+2] == 'l' {
				value := t.data[t.start : t.off+3]
				t.off += 3
				t.current = t.pendingNextStatus()
				return Null, value, nil
			}
			return Init, nil, fmt.Errorf("invalid null '%s'", string(t.data[t.start:]))
		}
	}

	switch t.current {
	case BeginObject,
		EndObject,
		BeginArray,
		EndArray,
		SepColon,
		SepComma:
		token := t.current
		t.current = EndJson
		return token, t.value, nil

	case EndJson:
		return EndJson, nil, nil

	case String:
		value := t.data[t.start:]
		t.current = EndJson
		if len(value) >= 2 && (value[len(value)-1] != '"' || value[len(value)-2] == '\\') {
			return String, value, nil
		}
		return Init, nil, fmt.Errorf("invalid string '%s'", string(value))

	case True:
		value := t.data[t.start:]
		t.current = EndJson
		if len(value) == 4 && value[1] == 'r' && value[2] == 'u' && value[3] == 'e' {
			return True, value, nil
		}
		return Init, nil, fmt.Errorf("invalid bool true '%s'", string(value))

	case False:
		value := t.data[t.start:]
		t.current = EndJson
		if len(value) == 5 && value[1] == 'a' && value[2] == 'l' && value[3] == 's' && value[4] == 'e' {
			return False, value, nil
		}
		return Init, nil, fmt.Errorf("invalid bool false '%s'", string(value))

	case Null:
		value := t.data[t.start:]
		t.current = EndJson
		if len(value) == 4 && value[1] == 'u' && value[2] == 'l' && value[3] == 'l' {
			return Null, value, nil
		}
		return Init, nil, fmt.Errorf("invalid null '%s'", string(value))

	default:
		token := t.current
		t.current = EndJson
		return token, t.data[t.start:], nil
	}
}
