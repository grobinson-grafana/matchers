package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLexer_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
		err      string
	}{{
		name:  "open paren",
		input: "{",
		expected: []Token{
			{Kind: TokenOpenParen, Value: "{"},
		},
	}, {
		name:  "open paren with space",
		input: " {",
		expected: []Token{
			{Kind: TokenOpenParen, Value: "{"},
		},
	}, {
		name:  "close paren",
		input: "}",
		expected: []Token{
			{Kind: TokenCloseParen, Value: "}"},
		},
	}, {
		name:  "close paren with space",
		input: "}",
		expected: []Token{
			{Kind: TokenCloseParen, Value: "}"},
		},
	}, {
		name:  "open and closing parens",
		input: "{}",
		expected: []Token{
			{Kind: TokenOpenParen, Value: "{"},
			{Kind: TokenCloseParen, Value: "}"},
		},
	}, {
		name:  "open and closing parens with space",
		input: "{ }",
		expected: []Token{
			{Kind: TokenOpenParen, Value: "{"},
			{Kind: TokenCloseParen, Value: "}"},
		},
	}, {
		name:  "ident",
		input: "hello",
		expected: []Token{
			{Kind: TokenIdent, Value: "hello"},
		},
	}, {
		name:  "idents with space",
		input: "hello world",
		expected: []Token{
			{Kind: TokenIdent, Value: "hello"},
			{Kind: TokenIdent, Value: "world"},
		},
	}, {
		name:  "quoted",
		input: "\"hello\"",
		expected: []Token{
			{Kind: TokenQuoted, Value: "\"hello\""},
		},
	}, {
		name:  "quoted with unicode",
		input: "\"hello 🙂\"",
		expected: []Token{
			{Kind: TokenQuoted, Value: "\"hello 🙂\""},
		},
	}, {
		name:  "quoted with space",
		input: "\"hello world\"",
		expected: []Token{
			{Kind: TokenQuoted, Value: "\"hello world\""},
		},
	}, {
		name:  "quoted with tab",
		input: "\"hello\tworld\"",
		expected: []Token{
			{Kind: TokenQuoted, Value: "\"hello\tworld\""},
		},
	}, {
		name:  "quoted with newline",
		input: "\"hello\nworld\"",
		expected: []Token{
			{Kind: TokenQuoted, Value: "\"hello\nworld\""},
		},
	}, {
		name:  "quoted with escaped quotes",
		input: "\"hello \\\"world\\\"\"",
		expected: []Token{
			{Kind: TokenQuoted, Value: "\"hello \\\"world\\\"\""},
		},
	}, {
		name:  "quoted with escaped backticks",
		input: "\"hello \\world\"",
		expected: []Token{
			{Kind: TokenQuoted, Value: "\"hello \\world\""},
		},
	}, {
		name:  "equals operator",
		input: "=",
		expected: []Token{
			{Kind: TokenOperator, Value: "="},
		},
	}, {
		name:  "not equals operator",
		input: "!=",
		expected: []Token{
			{Kind: TokenOperator, Value: "!="},
		},
	}, {
		name:  "matches regex operator",
		input: "=~",
		expected: []Token{
			{Kind: TokenOperator, Value: "=~"},
		},
	}, {
		name:  "not matches regex operator",
		input: "!~",
		expected: []Token{
			{Kind: TokenOperator, Value: "!~"},
		},
	}, {
		name:  "unexpected $",
		input: "$",
		err:   "unexpected input: $",
	}, {
		name:  "unexpected non alpha numeric in ident",
		input: "hello$",
		expected: []Token{
			{Kind: TokenIdent, Value: "hello"},
		},
		err: "unexpected input: hello$",
	}, {
		name:  "unexpected unicode in ident",
		input: "hello🙂",
		expected: []Token{
			{Kind: TokenIdent, Value: "hello"},
		},
		err: "unexpected input: hello🙂",
	}, {
		name:  "unexpected operator",
		input: "=$",
		expected: []Token{
			{Kind: TokenOperator, Value: "="},
		},
		err: "unexpected input: =$",
	}, {
		name:  "unterminated quoted",
		input: "\"hello",
		err:   "expected one of '\"', got EOF",
	}, {
		name:  "unterminated quoted with escaped quote",
		input: "\"hello\\\"",
		err:   "expected one of '\"', got EOF",
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := NewLexer(test.input)
			// scan all expected tokens
			for i := 0; i < len(test.expected); i++ {
				tok, err := l.Scan()
				require.NoError(t, err)
				assert.Equal(t, test.expected[i], tok)
			}
			if test.err == "" {
				// check there are no more tokens
				tok, err := l.Scan()
				require.NoError(t, err)
				assert.Equal(t, Token{}, tok)
			} else {
				// check if expected error is returned
				tok, err := l.Scan()
				assert.Equal(t, Token{}, tok)
				assert.EqualError(t, err, test.err)
			}
		})
	}
}

// This test asserts that the lexer does not emit more tokens after an
// error has occurred.
func TestLexer_ScanError(t *testing.T) {
	l := NewLexer("\"hello")
	for i := 0; i < 10; i++ {
		tok, err := l.Scan()
		assert.Equal(t, Token{}, tok)
		assert.EqualError(t, err, "expected one of '\"', got EOF")
	}
}
