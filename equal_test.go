package jsontools_test

import (
	"testing"

	"github.com/WqyJh/jsontools"
	"github.com/stretchr/testify/require"
)

func TestEqual(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{`{"a":null}`, `{}`, true},
		{`{"a":null,"b":"null"}`, `{"b":"null"}`, true},
		{`{"a":1 , "b":null}`, `{"a":1}`, true},
		{`{"a":null , "b":2}`, `{"b":2}`, true},

		{`{"a":1,"b":2,"c":null}`, `{"a":1,"b":2}`, true},
		{`{"a":1, "b":null,"c":3}`, `{"a":1,"c":3}`, true},
		{`{ "a":null,"b":2,"c":3}`, `{"b":2,"c":3}`, true},

		{`{"a":1,"b":null,"c":null}`, `{"a":1}`, true},
		{`{"a":null,"b":2, "c":null}`, `{"b":2}`, true},
		{`{"a":null,"b":null,"c":3}`, `{"c":3}`, true},

		{`{"a":null,"b":null,"c":null}`, `{}`, true},

		{`{"a":null,"b":null,"c":[1,2,3]}`, `{"c":[1,2,3]}`, true},
		{`{"a":null,"b":[1,2,3],"c":null}`, `{"b":[1,2,3]}`, true},
		{`{"a":[1,2,3] , "b" :null,"c":null }`, `{"a":[1,2,3]}`, true},
		{`{"a":null,"b":null,"c":[1,2,{ "a":1.23,"b":null ,"c"  :"null"}]}`, `{"c":[1,2,{"a":1.23,"c":"null"}]}`, true},
	}

	for _, c := range cases {
		got, err := jsontools.JsonEqual([]byte(c.a), []byte(c.b))
		require.NoError(t, err)
		require.Equal(t, c.want, got)
	}
}
