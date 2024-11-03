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
	parser := jsontools.NewJsonParser([]byte(expected1), func(ctx jsontools.HandlerContext) error {
		t.Logf("token: %v\tkind: %v\t\t'%s'", ctx.Token, ctx.Kind, string(ctx.Value))
		buf.Write(ctx.Value)
		return nil
	})
	require.NoError(t, parser.Parse())

	t.Logf("buf: '%s'", buf.String())
	require.JSONEq(t, expected1, buf.String())

	// stopped by error
	called := 0
	err := errors.New("stopped")
	parser = jsontools.NewJsonParser([]byte(expected1), func(ctx jsontools.HandlerContext) error {
		called++
		return err
	})
	require.ErrorIs(t, err, parser.Parse())
	require.Equal(t, 1, called)

	// with whitespace prefix/suffix
	expected := " \t\n" + expected1 + " \t\n"
	buf.Reset()
	parser = jsontools.NewJsonParser([]byte(expected), func(ctx jsontools.HandlerContext) error {
		t.Logf("token: %v\tkind: %v\t\t'%s'", ctx.Token, ctx.Kind, string(ctx.Value))
		buf.Write(ctx.Value)
		return nil
	})
	require.NoError(t, parser.Parse())

	t.Logf("buf: '%s'", buf.String())
	require.JSONEq(t, expected, buf.String())
}

func TestParser2(t *testing.T) {
	src := `[1, 2.3, true, false, null, "12345"]`
	expected := `[1,2.3,true,false,null,"12345"]`
	var buf bytes.Buffer
	parser := jsontools.NewJsonParser([]byte(src), func(ctx jsontools.HandlerContext) error {
		buf.Write(ctx.Value)
		return nil
	})
	require.NoError(t, parser.Parse())
	require.Equal(t, expected, buf.String())
}

func TestParserError(t *testing.T) {
	cases := []struct {
		src      string
		expected string
	}{
		{
			src:      ` `,
			expected: "invalid EOF1",
		},
		{
			src:      `{"hello"{`,
			expected: `invalid '{'`,
		},
		{
			src:      `"hello"`,
			expected: `invalid string '"hello"'`,
		},
		{
			src:      `12345`,
			expected: `invalid int '12345'`,
		},
		{
			src:      `{"key":"value" 1.23}`,
			expected: "invalid float '1.23'",
		},
		{
			src:      `{"key":"value", true}`,
			expected: "invalid bool 'true'",
		},
		{
			src:      `[]false`,
			expected: "invalid bool 'false'",
		},
		{
			src:      `{"key":"value"}null}`,
			expected: "invalid null 'null'",
		},
		{
			src:      `{"key":"value"},null}`,
			expected: "invalid ','",
		},
		{
			src:      `{"key"::"value"}}`,
			expected: "invalid ':'",
		},
		{
			src:      `{"key":"value"]}`,
			expected: "invalid ']'",
		},
		{
			src:      `[1,2}`,
			expected: "invalid '}'",
		},
		{
			src:      `{{`,
			expected: "invalid EOF2",
		},
	}
	for _, c := range cases {
		parser := jsontools.NewJsonParser([]byte(c.src), func(ctx jsontools.HandlerContext) error {
			return nil
		})
		err := parser.Parse()
		require.Error(t, err)
		require.Equal(t, c.expected, err.Error())
	}
}

func BenchmarkParser(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := jsontools.NewJsonParser([]byte(expected1), func(ctx jsontools.HandlerContext) error {
			return nil
		})
		require.NoError(b, parser.Parse())
	}
}
