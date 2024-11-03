package jsontools_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/WqyJh/jsontools"
	"github.com/stretchr/testify/require"
)

const (
	expected1 = `[
		{
		  "Field1": "12345",
		  "Field2": "MTIzNDU=",
		  "Field3": [1, 2.1, 3, -4.2, 5.00],
		  "Field4": {
			"1": "😄😄😄😄😄",
			"2": "12345",
			"3": "12345",
			"4": "12345",
			"5": "12345"
		  },
		  "Field5": {
			"Field1": "12345",
			"Field2": "MTIzNDU=",
			"Field3": [1, 2, 3, 4, 5],
			"Field4": {
			  "0": "12345",
			  "6": "12345",
			  "7": "12345",
			  "8": "12345",
			  "9": "12345"
			},
			"Field5": "12345",
			"Field6": { "Field1": "12345" },
			"Field7": ["12345", "12345", "12345", "12345", "12345"],
			"Field8": ["12345", "12345", "12345"],
			"Field9": ["12345", "12345", "12345", "12345", "12345", ""],
			"Field10": "12345"
		  },
		  "Field6": null,
		  "Field7": ["12345", "12345", "12345", "12345", "12345"],
		  "Field8": ["12345", "12345", "12345"],
		  "Field9": ["12345", "12345", "12345", "12345", "12345", ""],
		  "Field10": "12345"
		},
		{
		  "Field1": "12345",
		  "Field2": "MTIzNDU=",
		  "Field3": [1, 2, 3, 4, 5],
		  "Field4": {
			"0": "12345",
			"6": "12345",
			"7": "12345",
			"8": "12345",
			"9": "12345"
		  },
		  "Field5": "12345",
		  "Field6": { "Field1": "12345" },
		  "Field7": ["12345", "12345", "12345", "12345", "12345"],
		  "Field8": ["12345", "12345", "12345"],
		  "Field9": ["12345", "12345", "12345", "12345", "12345", ""],
		  "Field10": "12345"
		},
		{ "Field1": "11111" }
	  ]`
)

func TestTokenizer(t *testing.T) {
	require.True(t, json.Valid([]byte(expected1)))
	var buf bytes.Buffer
	tokenizer := jsontools.NewJsonTokenizer([]byte(expected1))
	for {
		token, value, err := tokenizer.Next()
		if err != nil {
			break
		}
		t.Logf("token: %v\t\t'%s'", token, string(value))
		if token == jsontools.EndJson {
			break
		}
		buf.Write(value)
	}
	t.Logf("buf: '%s'", buf.String())
	require.JSONEq(t, expected1, buf.String())
}

func TestTokenizerError(t *testing.T) {
	cases := []struct {
		src      string
		expected string
	}{
		{`1..`, "invalid float"},
		{`t`, "invalid bool true 't'"},
		{`tr`, "invalid bool true 'tr'"},
		{`tru`, "invalid bool true 'tru'"},
		{`trus`, "invalid bool true 'trus'"},
		{`talse`, "invalid bool true 'talse'"},
		{`f`, "invalid bool false 'f'"},
		{`fa`, "invalid bool false 'fa'"},
		{`fal`, "invalid bool false 'fal'"},
		{`fals`, "invalid bool false 'fals'"},
		{`falss`, "invalid bool false 'falss'"},
		{`falsse`, "invalid bool false 'falsse'"},
		{`n`, "invalid null 'n'"},
		{`nu`, "invalid null 'nu'"},
		{`nul`, "invalid null 'nul'"},
		{`nuls`, "invalid null 'nuls'"},
		{`nulsl`, "invalid null 'nulsl'"},
		{`"`, `invalid string '"'`},
		{`"1`, `invalid string '"1'`},
		{`"1\"`, `invalid string '"1\"'`},
	}
	for _, c := range cases {
		tokenizer := jsontools.NewJsonTokenizer([]byte(c.src))
		t.Logf("src: '%s'", c.src)
		_, _, err := tokenizer.Next()
		require.Error(t, err)
		require.Equal(t, c.expected, err.Error())
	}
}

func BenchmarkTokenizer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokenizer := jsontools.NewJsonTokenizer([]byte(expected1))
		for {
			token, _, err := tokenizer.Next()
			require.NoError(b, err)
			if token == jsontools.EndJson {
				break
			}
		}
	}
}
