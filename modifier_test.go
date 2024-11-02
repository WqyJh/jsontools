package jsontools_test

import (
	"bytes"
	"testing"

	"github.com/WqyJh/jsontools"
	"github.com/stretchr/testify/require"
)

func TestModifyJson(t *testing.T) {
	src := `[{"Field1":"1234567890","Field2":"MTIzNDU2Nzg5MA==","Field3":[1,2,3,4,5,6,7,8,9,0],"Field4":{"1":"1234567890","2":"1234567890","3":"1234567890","4":"1234567890","5":"1234567890","6":"1234567890"},"Field5":{"Field1":"1234567890","Field2":"MTIzNDU2Nzg5MA==","Field3":[1,2,3,4,5,6,7,8,9,0],"Field4":{"0":"1234567890","6":"1234567890","7":"1234567890","8":"1234567890","9":"1234567890"},"Field5":"1234567890","Field6":{"Field1":"1234567890"},"Field7":["1234567890","1234567890","1234567890","1234567890","1234567890","1234567890"],"Field8":["1234567890","1234567890","1234567890"],"Field9":["1234567890","1234567890","1234567890","1234567890","1234567890","1234567890"],"Field10":"1234567890"},"Field6":null,"Field7":["1234567890","1234567890","1234567890","1234567890","1234567890","1234567890"],"Field8":["1234567890","1234567890","1234567890"],"Field9":["1234567890","1234567890","1234567890","1234567890","1234567890","1234567890"],"Field10":"1234567890"},{"Field1":"1234567890","Field2":"MTIzNDU2Nzg5MA==","Field3":[1,2,3,4,5,6,7,8,9,0],"Field4":{"0":"1234567890","6":"1234567890","7":"1234567890","8":"1234567890","9":"1234567890"},"Field5":"1234567890","Field6":{"Field1":"1234567890"},"Field7":["1234567890","1234567890","1234567890","1234567890","1234567890","1234567890"],"Field8":["1234567890","1234567890","1234567890"],"Field9":["1234567890","1234567890","1234567890","1234567890","1234567890","1234567890"],"Field10":"1234567890"},{"Field1":"1234567890"}]`
	expected := `[{"Field1":"12345","Field2":"MTIzN","Field3":[1,2,3,4,5,6,7,8,9,0],"Field4":{"1":"12345","2":"12345","3":"12345","4":"12345","5":"12345","6":"12345"},"Field5":{"Field1":"12345","Field2":"MTIzN","Field3":[1,2,3,4,5,6,7,8,9,0],"Field4":{"0":"12345","6":"12345","7":"12345","8":"12345","9":"12345"},"Field5":"12345","Field6":{"Field1":"12345"},"Field7":["12345","12345","12345","12345","12345","12345"],"Field8":["12345","12345","12345"],"Field9":["12345","12345","12345","12345","12345","12345"],"Field10":"12345"},"Field6":null,"Field7":["12345","12345","12345","12345","12345","12345"],"Field8":["12345","12345","12345"],"Field9":["12345","12345","12345","12345","12345","12345"],"Field10":"12345"},{"Field1":"12345","Field2":"MTIzN","Field3":[1,2,3,4,5,6,7,8,9,0],"Field4":{"0":"12345","6":"12345","7":"12345","8":"12345","9":"12345"},"Field5":"12345","Field6":{"Field1":"12345"},"Field7":["12345","12345","12345","12345","12345","12345"],"Field8":["12345","12345","12345"],"Field9":["12345","12345","12345","12345","12345","12345"],"Field10":"12345"},{"Field1":"12345"}]`
	srcBytes := bytes.Clone([]byte(src))

	// dstBytes is a new bytes, not change srcBytes
	dstBytes, err := jsontools.ModifyJson(srcBytes, jsontools.WithFieldLengthLimit(5))
	require.NoError(t, err)
	require.JSONEq(t, expected, string(dstBytes))
	require.Equal(t, src, string(srcBytes))

	// srcBytes would be changed
	dstBytes, err = jsontools.ModifyJson(srcBytes, jsontools.WithFieldLengthLimit(5), jsontools.WithInplace(true))
	require.NoError(t, err)
	require.JSONEq(t, expected, string(dstBytes))
	require.NotEqual(t, expected, string(srcBytes))
	// same prefix
	require.Equal(t, expected, string(srcBytes[:len(dstBytes)]))

	// check all string length
	parser := jsontools.NewJsonParser([]byte(expected), func(token jsontools.TokenType, kind jsontools.Kind, value []byte) {
		switch kind {
		case jsontools.KindArrayValue, jsontools.KindObjectValue:
			if token == jsontools.String {
				require.LessOrEqual(t, len(value), 5+2)
			}
		}
	})
	require.NoError(t, parser.Parse())

	// not change because all string length is less than limit
	srcBytes = bytes.Clone([]byte(src))
	dstBytes, err = jsontools.ModifyJson(srcBytes, jsontools.WithFieldLengthLimit(100), jsontools.WithInplace(true))
	require.NoError(t, err)
	require.Equal(t, src, string(dstBytes))
	require.Equal(t, src, string(srcBytes))
}
