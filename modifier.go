package jsontools

type jsonModifier struct {
	limit   int
	inplace bool
}

type jsonModifierOption func(*jsonModifier)

func WithFieldLengthLimit(limit int) jsonModifierOption {
	return func(m *jsonModifier) {
		m.limit = limit + 2
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
	parser := NewJsonParser(data, func(token TokenType, kind Kind, value []byte) {
		needModify := false
		if modifier.limit > 0 {
			switch kind {
			case KindObjectValue,
				KindArrayValue:
				if token == String && len(value) > modifier.limit {
					needModify = true
				}
			}
		}
		if needModify {
			dst = append(dst, value[:modifier.limit-1]...)
			dst = append(dst, '"')
		} else {
			dst = append(dst, value...)
		}
	})
	err := parser.Parse()
	if err != nil {
		return nil, err
	}
	return dst, nil
}
