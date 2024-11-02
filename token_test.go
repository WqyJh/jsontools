package jsontools_test

import (
	"bytes"
	"testing"

	"github.com/WqyJh/jsontools"
	"github.com/stretchr/testify/require"
)

func TestTokenizer(t *testing.T) {
	expected := `[
		{
		  "Field1": "12345",
		  "Field2": "MTIzNDU=",
		  "Field3": [1, 2, 3, 4, 5],
		  "Field4": {
			"1": "12345",
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
	var buf bytes.Buffer
	tokenizer := jsontools.NewJsonTokenizer([]byte(expected))
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
	require.JSONEq(t, expected, buf.String())
}
