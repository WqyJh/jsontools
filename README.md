# jsontools [![GoDoc][doc-img]][doc]


Simple tools for json in golang.

Features:
- Tokenize json bytes.
- Parse and validate json bytes.
- Modify json string with field length limit.
- Filter null values from json bytes.
- Check if two json bytes are equal except null values.

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

    parser := jsontools.NewJsonParser([]byte(expected), func(ctx jsontools.HandlerContext) error {
		fmt.Printf("token: %v\tkind: %v\t\t'%s'\n", ctx.Token, ctx.Kind, string(ctx.Value))
		return nil
	})
```

Supported kinds:

| Kind         | Representation |
|--------------|----------------|
| ObjectKey    | object key     |
| ObjectValue  | object value   |
| ArrayValue   | array value    |

If you return error in handler, the parser will be stopped.


### Modify Json

When you got a json string, and you want to write it to log, but some of the fields are too long, you can use this tool to modify the json string, by cutting off the long fields.

There are `inplace` mode to modify the input bytes directly, without allocating new bytes, used only when src won't be used anymore.

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

Or if you want to filter some keys from the output, such as password or credentials, use the following.

```go
src := `{"a":"1234567890","b":"1234567890","c":"1234567890","d":"1234567890"}`

// result is `{"a":"12345","c":"12345"}`
dst, err = jsontools.ModifyJson([]byte(src), jsontools.WithFilterKeys("b", "d"), jsontools.WithFieldLengthLimit(5), jsontools.WithInplace(true))
```

`ModifyJson` is a wrapper of `JsonModifier`, which create a new `JsonModifier` on every call. If you want to modify multiple json strings with same options, you can create a `JsonModifier` once, and call `JsonModifier.ModifyJson` method multiple times, which is a concurrent-safe reentrant function.

```go
modifier := jsontools.NewJsonModifier(jsontools.WithFilterKeys("b", "d"), jsontools.WithFieldLengthLimit(5))

// result is `{"a":"12345","c":"12345"}`
dst, err = modifier.ModifyJson([]byte(src))
```

### Filter Null

Filter null values from json bytes.

```go
src := `{"a":"1234567890","b":null,"c":"null"}`

filter := jsontools.NewJsonNullFilter(false)

// result is `{"a":"12345","c":"null"}`
dst, err = filter.Filter([]byte(src))
```

### Json Equal

Check if two json bytes are equal except null values.

```go
src1 := `{"a":"1234567890","b":null,"c":"null"}`
src2 := `{"a":"12345","c":"null"}`

// equal is true
equal, err := jsontools.JsonEqual([]byte(src1), []byte(src2))
```

Sometimes we want to check if `json.Marshal` result is expected by using `assert.JSONEq` from [github.com/stretchr/testify/assert](https://pkg.go.dev/github.com/stretchr/testify/assert#JSONEq). But if some fields are `omitempty`, the marshal result won't contain these fields, however the expected json string may contain null values of these fields, which cause `assert.JSONEq` failed. Use `JsonEqual` to check if two json bytes are equal except null values, which is useful in this case.


## License

Released under the [MIT License](LICENSE).

[doc-img]: https://godoc.org/github.com/WqyJh/jsontools?status.svg
[doc]: https://godoc.org/github.com/WqyJh/jsontools
