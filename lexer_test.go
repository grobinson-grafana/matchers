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
		expected: []Token{{
			Kind:  TokenOpenParen,
			Value: "{",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   1,
				ColumnStart: 0,
				ColumnEnd:   1,
			},
		}},
	}, {
		name:  "open paren with space",
		input: " {",
		expected: []Token{{
			Kind:  TokenOpenParen,
			Value: "{",
			Position: Position{
				OffsetStart: 1,
				OffsetEnd:   2,
				ColumnStart: 1,
				ColumnEnd:   2,
			},
		}},
	}, {
		name:  "close paren",
		input: "}",
		expected: []Token{{
			Kind:  TokenCloseParen,
			Value: "}",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   1,
				ColumnStart: 0,
				ColumnEnd:   1,
			},
		}},
	}, {
		name:  "close paren with space",
		input: "}",
		expected: []Token{{
			Kind:  TokenCloseParen,
			Value: "}",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   1,
				ColumnStart: 0,
				ColumnEnd:   1,
			},
		}},
	}, {
		name:  "open and closing parens",
		input: "{}",
		expected: []Token{{
			Kind:  TokenOpenParen,
			Value: "{",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   1,
				ColumnStart: 0,
				ColumnEnd:   1,
			},
		}, {
			Kind:  TokenCloseParen,
			Value: "}",
			Position: Position{
				OffsetStart: 1,
				OffsetEnd:   2,
				ColumnStart: 1,
				ColumnEnd:   2,
			},
		}},
	}, {
		name:  "open and closing parens with space",
		input: "{ }",
		expected: []Token{{
			Kind:  TokenOpenParen,
			Value: "{",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   1,
				ColumnStart: 0,
				ColumnEnd:   1,
			},
		}, {
			Kind:  TokenCloseParen,
			Value: "}",
			Position: Position{
				OffsetStart: 2,
				OffsetEnd:   3,
				ColumnStart: 2,
				ColumnEnd:   3,
			},
		}},
	}, {
		name:  "ident",
		input: "hello",
		expected: []Token{{
			Kind:  TokenIdent,
			Value: "hello",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   5,
				ColumnStart: 0,
				ColumnEnd:   5,
			},
		}},
	}, {
		name:  "ident with underscore",
		input: "hello_world",
		expected: []Token{{
			Kind:  TokenIdent,
			Value: "hello_world",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   11,
				ColumnStart: 0,
				ColumnEnd:   11,
			},
		}},
	}, {
		name:  "ident with colon",
		input: "hello:world",
		expected: []Token{{
			Kind:  TokenIdent,
			Value: "hello:world",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   11,
				ColumnStart: 0,
				ColumnEnd:   11,
			},
		}},
	}, {
		name:  "ident with numbers",
		input: "hello0123456789",
		expected: []Token{{
			Kind:  TokenIdent,
			Value: "hello0123456789",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   15,
				ColumnStart: 0,
				ColumnEnd:   15,
			},
		}},
	}, {
		name:  "ident can start with underscore",
		input: "_hello",
		expected: []Token{{
			Kind:  TokenIdent,
			Value: "_hello",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   6,
				ColumnStart: 0,
				ColumnEnd:   6,
			},
		}},
	}, {
		name:  "idents separated with space",
		input: "hello world",
		expected: []Token{{
			Kind:  TokenIdent,
			Value: "hello",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   5,
				ColumnStart: 0,
				ColumnEnd:   5,
			},
		}, {
			Kind:  TokenIdent,
			Value: "world",
			Position: Position{
				OffsetStart: 6,
				OffsetEnd:   11,
				ColumnStart: 6,
				ColumnEnd:   11,
			},
		}},
	}, {
		name:  "quoted",
		input: "\"hello\"",
		expected: []Token{{
			Kind:  TokenQuoted,
			Value: "\"hello\"",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   7,
				ColumnStart: 0,
				ColumnEnd:   7,
			},
		}},
	}, {
		name:  "quoted with unicode",
		input: "\"hello 🙂\"",
		expected: []Token{{
			Kind:  TokenQuoted,
			Value: "\"hello 🙂\"",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   12,
				ColumnStart: 0,
				ColumnEnd:   9,
			},
		}},
	}, {
		name:  "quoted with space",
		input: "\"hello world\"",
		expected: []Token{{
			Kind:  TokenQuoted,
			Value: "\"hello world\"",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   13,
				ColumnStart: 0,
				ColumnEnd:   13,
			},
		}},
	}, {
		name:  "quoted with newline",
		input: "\"hello\nworld\"",
		expected: []Token{{
			Kind:  TokenQuoted,
			Value: "\"hello\nworld\"",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   13,
				ColumnStart: 0,
				ColumnEnd:   13,
			},
		}},
	}, {
		name:  "quoted with tab",
		input: "\"hello\tworld\"",
		expected: []Token{{
			Kind:  TokenQuoted,
			Value: "\"hello\tworld\"",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   13,
				ColumnStart: 0,
				ColumnEnd:   13,
			},
		}},
	}, {
		name:  "quoted with escaped quotes",
		input: "\"hello \\\"world\\\"\"",
		expected: []Token{{
			Kind:  TokenQuoted,
			Value: "\"hello \\\"world\\\"\"",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   17,
				ColumnStart: 0,
				ColumnEnd:   17,
			},
		}},
	}, {
		name:  "quoted with escaped backslash",
		input: "\"hello world\\\\\"",
		expected: []Token{{
			Kind:  TokenQuoted,
			Value: "\"hello world\\\\\"",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   15,
				ColumnStart: 0,
				ColumnEnd:   15,
			},
		}},
	}, {
		name:  "equals operator",
		input: "=",
		expected: []Token{{
			Kind:  TokenOperator,
			Value: "=",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   1,
				ColumnStart: 0,
				ColumnEnd:   1,
			},
		}},
	}, {
		name:  "not equals operator",
		input: "!=",
		expected: []Token{{
			Kind:  TokenOperator,
			Value: "!=",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   2,
				ColumnStart: 0,
				ColumnEnd:   2,
			},
		}},
	}, {
		name:  "matches regex operator",
		input: "=~",
		expected: []Token{{
			Kind:  TokenOperator,
			Value: "=~",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   2,
				ColumnStart: 0,
				ColumnEnd:   2,
			},
		}},
	}, {
		name:  "not matches regex operator",
		input: "!~",
		expected: []Token{{
			Kind:  TokenOperator,
			Value: "!~",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   2,
				ColumnStart: 0,
				ColumnEnd:   2,
			},
		}},
	}, {
		name:  "unexpected $",
		input: "$",
		err:   "0:1: $: invalid input",
	}, {
		name:  "unexpected emoji",
		input: "🙂",
		err:   "0:1: 🙂: invalid input",
	}, {
		name:  "unexpected unicode letter",
		input: "Σ",
		err:   "0:1: Σ: invalid input",
	}, {
		name:  "unexpected : at start of ident",
		input: ":hello",
		err:   "0:1: :: invalid input",
	}, {
		name:  "unexpected $ in ident",
		input: "hello$",
		expected: []Token{{
			Kind:  TokenIdent,
			Value: "hello",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   5,
				ColumnStart: 0,
				ColumnEnd:   5,
			},
		}},
		err: "5:6: $: invalid input",
	}, {
		name:  "unexpected unicode letter in ident",
		input: "helloΣ",
		expected: []Token{{
			Kind:  TokenIdent,
			Value: "hello",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   5,
				ColumnStart: 0,
				ColumnEnd:   5,
			},
		}},
		err: "5:6: Σ: invalid input",
	}, {
		name:  "unexpected emoji in ident",
		input: "hello🙂",
		expected: []Token{{
			Kind:  TokenIdent,
			Value: "hello",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   5,
				ColumnStart: 0,
				ColumnEnd:   5,
			},
		}},
		err: "5:6: 🙂: invalid input",
	}, {
		name:  "invalid operator",
		input: "!",
		err:   "0:1: unexpected end of input, expected one of '=~'",
	}, {
		name:  "another invalid operator",
		input: "~",
		err:   "0:1: ~: invalid input",
	}, {
		name:  "unexpected $ in operator",
		input: "=$",
		expected: []Token{{
			Kind:  TokenOperator,
			Value: "=",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   1,
				ColumnStart: 0,
				ColumnEnd:   1,
			},
		}},
		err: "1:2: $: invalid input",
	}, {
		name:  "unexpected ! after operator",
		input: "=!",
		expected: []Token{{
			Kind:  TokenOperator,
			Value: "=",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   1,
				ColumnStart: 0,
				ColumnEnd:   1,
			},
		}},
		err: "1:2: unexpected end of input, expected one of '=~'",
	}, {
		name:  "unexpected !! after operator",
		input: "!=!!",
		expected: []Token{{
			Kind:  TokenOperator,
			Value: "!=",
			Position: Position{
				OffsetStart: 0,
				OffsetEnd:   2,
				ColumnStart: 0,
				ColumnEnd:   2,
			},
		}},
		err: "2:3: !: expected one of '=~'",
	}, {
		name:  "unterminated quoted",
		input: "\"hello",
		err:   "0:6: \"hello: missing end \"",
	}, {
		name:  "unterminated quoted with escaped quote",
		input: "\"hello\\\"",
		err:   "0:8: \"hello\\\": missing end \"",
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
		assert.EqualError(t, err, "0:6: \"hello: missing end \"")
	}
}

func TestLexer_Peek(t *testing.T) {
	l := NewLexer("hello world")
	expected1 := Token{
		Kind:  TokenIdent,
		Value: "hello",
		Position: Position{
			OffsetStart: 0,
			OffsetEnd:   5,
			ColumnStart: 0,
			ColumnEnd:   5,
		},
	}
	expected2 := Token{
		Kind:  TokenIdent,
		Value: "world",
		Position: Position{
			OffsetStart: 6,
			OffsetEnd:   11,
			ColumnStart: 6,
			ColumnEnd:   11,
		},
	}
	// check that Peek() returns the first token
	tok, err := l.Peek()
	assert.NoError(t, err)
	assert.Equal(t, expected1, tok)
	// check that Scan() returns the peeked token
	tok, err = l.Scan()
	assert.NoError(t, err)
	assert.Equal(t, expected1, tok)
	// check that Peek() returns the second token until the next Scan()
	for i := 0; i < 10; i++ {
		tok, err = l.Peek()
		assert.NoError(t, err)
		assert.Equal(t, expected2, tok)
	}
	// check that Scan() returns the last token
	tok, err = l.Scan()
	assert.NoError(t, err)
	assert.Equal(t, expected2, tok)
	// should not be able to Peek() further tokens
	for i := 0; i < 10; i++ {
		tok, err = l.Peek()
		assert.NoError(t, err)
		assert.Equal(t, Token{}, tok)
	}
}

// This test asserts that the lexer does not emit more tokens after an
// error has occurred.
func TestLexer_PeekError(t *testing.T) {
	l := NewLexer("\"hello")
	for i := 0; i < 10; i++ {
		tok, err := l.Peek()
		assert.Equal(t, Token{}, tok)
		assert.EqualError(t, err, "0:6: \"hello: missing end \"")
	}
}
