# jsontools [![GoDoc][doc-img]][doc]


Simple tools for json in golang.

Features:
- Tokenize json bytes.
- Parse and validate json bytes.
- Modify json string with field length limit.

## Get

```bash
go get -u "github.com/WqyJh/jsontools"
```

## Usage

### Tokenizer

Iterate all tokens of input json bytes. 
```go
import (
	"github.com/WqyJh/jsontools"
)

    tokenizer := jsontools.NewJsonTokenizer([]byte(expected))
	for {
		token, value, err := tokenizer.Next()
		if err != nil {
			break
		}
		if token == jsontools.EndJson {
			break
		}
        // do something with value
		fmt.Printf("token: %v\t\t'%s'\n", token, string(value))
	}
```

Supported token types:

| Token Type   | Representation |
|--------------|----------------|
| BeginObject  | {              |
| EndObject    | }              |
| BeginArray   | [              |
| EndArray     | ]              |
| Null         | null           |
| Number       | number         |
| Float        | float          |
| String       | "string"      |
| True         | true           |
| False        | false          |
| SepColon     | :              |
| SepComma     | ,              |


### Parser

The tokenizer is used to analyze json bytes by syntax, but it's not enough for some cases.

For example, if you get a string token, and you want to know whether it's an object key, an object value or an array value, you need to use parser.

```go
import (
	"github.com/WqyJh/jsontools"
)

    parser := jsontools.NewJsonParser([]byte(expected), func(token jsontools.TokenType, kind jsontools.Kind, value []byte) {
		fmt.Printf("token: %v\tkind: %v\t\t'%s'\n", token, kind, string(value))
	})
```

Supported kinds:

| Kind         | Representation |
|--------------|----------------|
| ObjectKey    | object key     |
| ObjectValue  | object value   |
| ArrayValue   | array value    |


### Modify Json

When you got a json string, and you want to write it to log, but some of the fields are too long, you can use this tool to modify the json string, by cutting off the long fields.

```go
import (
	"github.com/WqyJh/jsontools"
)

src := `{"a":"1234567890","b":"1234567890","c":"1234567890"}`

// result is `{"a":"12345","b":"12345","c":"12345"}`
dst, _ := jsontools.ModifyJson([]byte(src), jsontools.WithFieldLengthLimit(5))

// result is `{"a":"12345","b":"12345","c":"12345"}`
// and src is modified, used only when src won't be used anymore
dst, _ = jsontools.ModifyJson([]byte(src), jsontools.WithFieldLengthLimit(5), jsontools.WithInplace(true))
```

There are `inplace` mode to modify the input bytes directly, without allocating new bytes, used only when src won't be used anymore.

## License

Released under the [MIT License](LICENSE).

[doc-img]: https://godoc.org/github.com/WqyJh/jsontools?status.svg
[doc]: https://godoc.org/github.com/WqyJh/jsontools
