package jsontools_test

import (
	"bytes"
	"testing"

	"github.com/WqyJh/jsontools"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	var buf bytes.Buffer
	parser := jsontools.NewJsonParser([]byte(expected1), func(token jsontools.TokenType, kind jsontools.Kind, value []byte) {
		t.Logf("token: %v\tkind: %v\t\t'%s'", token, kind, string(value))
		buf.Write(value)
	})
	require.NoError(t, parser.Parse())

	t.Logf("buf: '%s'", buf.String())
	require.JSONEq(t, expected1, buf.String())
}

func BenchmarkParser(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := jsontools.NewJsonParser([]byte(expected1), func(token jsontools.TokenType, kind jsontools.Kind, value []byte) {})
		require.NoError(b, parser.Parse())
	}
}
