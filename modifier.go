package jsontools

import (
	"errors"
	"unicode/utf8"
)

type jsonModifier struct {
	limit   int
	inplace bool
}

type jsonModifierOption func(*jsonModifier)

func WithFieldLengthLimit(limit int) jsonModifierOption {
	return func(m *jsonModifier) {
		m.limit = limit
	}
}

func WithInplace(inplace bool) jsonModifierOption {
	return func(m *jsonModifier) {
		m.inplace = inplace
	}
}

func ModifyJson(data []byte, opts ...jsonModifierOption) ([]byte, error) {
	modifier := &jsonModifier{}
	for _, opt := range opts {
		opt(modifier)
	}
	var dst []byte
	if modifier.inplace {
		dst = data[:0]
	} else {
		dst = make([]byte, 0, len(data))
	}
	parser := NewJsonParser(data, func(token TokenType, kind Kind, value []byte) error {
		needModify := false
		if modifier.limit > 0 {
			switch kind {
			case KindObjectValue,
				KindArrayValue:
				if token == String && utf8.RuneCount(value) > modifier.limit+2 {
					needModify = true
				}
			}
		}
		if needModify {
			dst = append(dst, '"')
			count := 0
			for i := 1; ; {
				r, size := utf8.DecodeRune(value[i:])
				if r == utf8.RuneError {
					return errors.New("invalid utf8")
				}
				dst = append(dst, value[i:i+size]...)
				i += size
				count++
				if count >= modifier.limit {
					break
				}
			}
			dst = append(dst, '"')
		} else {
			dst = append(dst, value...)
		}
		return nil
	})
	err := parser.Parse()
	if err != nil {
		return nil, err
	}
	return dst, nil
}
