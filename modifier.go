package jsontools

import (
	"errors"
	"unicode/utf8"
)

type jsonModifier struct {
	limit   int
	inplace bool

	filterKeySet    map[string]struct{}
	skipComma       bool
	expectStackSize int
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

func WithFilterKeys(keys ...string) jsonModifierOption {
	return func(m *jsonModifier) {
		m.filterKeySet = make(map[string]struct{}, len(keys))
		for _, key := range keys {
			k := `"` + key + `"`
			m.filterKeySet[k] = struct{}{}
		}
	}
}

func ModifyJson(data []byte, opts ...jsonModifierOption) ([]byte, error) {
	m := &jsonModifier{}
	for _, opt := range opts {
		opt(m)
	}
	var dst []byte
	if m.inplace {
		dst = data[:0]
	} else {
		dst = make([]byte, 0, len(data))
	}

	filterKeyStack := make([][]byte, 0, 32)
	parser := NewJsonParser(data, func(ctx HandlerContext) error {

		// filter keys ------- begin -------
		if m.expectStackSize > 0 {
			if ctx.StackSize >= m.expectStackSize {
				// skip all colon, comma and values of this key
				return nil
			} else {
				m.expectStackSize = 0
				m.skipComma = true
			}
		}

		switch ctx.Kind {
		case KindObjectKey:
			// object key must be string, therefore, this if could be removed
			if _, ok := m.filterKeySet[string(ctx.Value)]; ok {
				filterKeyStack = append(filterKeyStack, ctx.Value)
				// skip this key
				// m.skipColon = true
				m.expectStackSize = ctx.StackSize
				return nil
			}

		case KindOther:
			if m.skipComma {
				m.skipComma = false
				if ctx.Token == SepComma {
					// skip this comma
					return nil
				}
			}
		}
		// filter keys ------- end -------

		// modify value ------- begin -------
		needModify := false
		if m.limit > 0 {
			switch ctx.Kind {
			case KindObjectValue,
				KindArrayValue:
				if ctx.Token == String && utf8.RuneCount(ctx.Value) > m.limit+2 {
					needModify = true
				}
			}
		}
		if needModify {
			dst = append(dst, '"')
			count := 0
			for i := 1; ; {
				r, size := utf8.DecodeRune(ctx.Value[i:])
				if r == utf8.RuneError {
					return errors.New("invalid utf8")
				}
				dst = append(dst, ctx.Value[i:i+size]...)
				i += size
				count++
				if count >= m.limit {
					break
				}
			}
			dst = append(dst, '"')
		} else {
			dst = append(dst, ctx.Value...)
		}
		// modify value ------- end -------

		return nil
	})
	err := parser.Parse()
	if err != nil {
		return nil, err
	}
	return dst, nil
}
