package jsontools

type JsonNullFilter struct {
	inplace bool
}

func NewJsonNullFilter(inplace bool) *JsonNullFilter {
	return &JsonNullFilter{inplace: inplace}
}

func (f *JsonNullFilter) Filter(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	var dst []byte
	if f.inplace {
		dst = data[:0]
	} else {
		dst = make([]byte, 0, len(data))
	}

	lastPos := -1
	skipComma := false
	commaSkipped := false
	parser := NewJsonParser(data, func(ctx HandlerContext) error {

		switch ctx.Kind {
		case KindObjectKey:
			if commaSkipped {
				if dst[len(dst)-1] != '{' {
					// restore the comma for next key
					dst = append(dst, ',')
				}
				commaSkipped = false
			}
		default:
			commaSkipped = false
		}

		switch ctx.Kind {
		case KindObjectKey:

		case KindObjectValue:
			if lastPos != -1 {
				switch ctx.Token {
				case Null:
					dst = dst[:lastPos] // remove the key
					skipComma = true
					lastPos = -1
					return nil
				default:
					lastPos = -1
				}
			}

		case KindArrayValue:
			if lastPos != -1 {
				lastPos = -1
			}

		case KindOther:
			switch ctx.Token {
			case SepComma:
				lastPos = len(dst)
			}

			if skipComma {
				skipComma = false
				if ctx.Token == SepComma {
					commaSkipped = true
					// skip this comma
					return nil
				}
			}
		}

		dst = append(dst, ctx.Value...)

		switch ctx.Token {
		case BeginObject:
			lastPos = len(dst)
		}
		return nil
	})

	err := parser.Parse()
	if err != nil {
		return nil, err
	}
	return dst, nil
}
