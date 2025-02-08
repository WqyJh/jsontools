package jsontools_test

import (
	"testing"

	"github.com/WqyJh/jsontools"
	"github.com/stretchr/testify/require"
)

func TestJsonNullFilter(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{`{"a":null}`, `{}`},
		{`{"a":null,"b":"null"}`, `{"b":"null"}`},

		{`{"a":1 , "b":null}`, `{"a":1}`},
		{`{"a":null , "b":2}`, `{"b":2}`},

		{`{"a":1,"b":2,"c":null}`, `{"a":1,"b":2}`},
		{`{"a":1, "b":null,"c":3}`, `{"a":1,"c":3}`},
		{`{ "a":null,"b":2,"c":3}`, `{"b":2,"c":3}`},

		{`{"a":1,"b":null,"c":null}`, `{"a":1}`},
		{`{"a":null,"b":2, "c":null}`, `{"b":2}`},
		{`{"a":null,"b":null,"c":3}`, `{"c":3}`},

		{`{"a":null,"b":null,"c":null}`, `{}`},

		{`{"a":null,"b":null,"c":[1,2,3]}`, `{"c":[1,2,3]}`},
		{`{"a":null,"b":[1,2,3],"c":null}`, `{"b":[1,2,3]}`},
		{`{"a":[1,2,3] , "b" :null,"c":null }`, `{"a":[1,2,3]}`},
		{`{"a":null,"b":null,"c":[1,2,{ "a":1.23,"b":null ,"c"  :"null"}]}`, `{"c":[1,2,{"a":1.23,"c":"null"}]}`},
	}
	for i, c := range cases {
		dst, err := jsontools.NewJsonNullFilter(false).Filter([]byte(c.input))
		require.NoError(t, err)
		require.Equal(t, c.expected, string(dst), "case %d: %s", i, c.input)

		// inplace
		input := []byte(c.input) // cloned
		dst, err = jsontools.NewJsonNullFilter(true).Filter(input)
		require.NoError(t, err)
		require.Equal(t, c.expected, string(dst), "case %d: dst:%s src:%s", i, dst, input)
		t.Logf("inplace case %d: dst:%s src:%s", i, dst, input)
	}
}
