package jsontools_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/WqyJh/jsontools"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	var buf bytes.Buffer
	parser := jsontools.NewJsonParser([]byte(expected1), func(token jsontools.TokenType, kind jsontools.Kind, value []byte) error {
		t.Logf("token: %v\tkind: %v\t\t'%s'", token, kind, string(value))
		buf.Write(value)
		return nil
	})
	require.NoError(t, parser.Parse())

	t.Logf("buf: '%s'", buf.String())
	require.JSONEq(t, expected1, buf.String())

	// stopped by error
	called := 0
	err := errors.New("stopped")
	parser = jsontools.NewJsonParser([]byte(expected1), func(token jsontools.TokenType, kind jsontools.Kind, value []byte) error {
		called++
		return err
	})
	require.ErrorIs(t, err, parser.Parse())
	require.Equal(t, 1, called)
}

func BenchmarkParser(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := jsontools.NewJsonParser([]byte(expected1), func(token jsontools.TokenType, kind jsontools.Kind, value []byte) error {
			return nil
		})
		require.NoError(b, parser.Parse())
	}
}
