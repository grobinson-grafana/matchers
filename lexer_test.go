package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{{
		name:  "open and closing parens",
		input: "{}",
		expected: []Token{
			{Kind: TokenOpenParen, Value: "{"},
			{Kind: TokenCloseParen, Value: "}"},
		},
	}, {
		name:  "equals",
		input: "{hello=world}",
		expected: []Token{
			{Kind: TokenOpenParen, Value: "{"},
			{Kind: TokenLiteral, Value: "hello"},
			{Kind: TokenOperator, Value: "="},
			{Kind: TokenLiteral, Value: "world"},
			{Kind: TokenCloseParen, Value: "}"},
		},
	}, {
		name:  "not equals",
		input: "{hello!=world}",
		expected: []Token{
			{Kind: TokenOpenParen, Value: "{"},
			{Kind: TokenLiteral, Value: "hello"},
			{Kind: TokenOperator, Value: "!="},
			{Kind: TokenLiteral, Value: "world"},
			{Kind: TokenCloseParen, Value: "}"},
		},
	}, {
		name:  "match regex",
		input: "{hello=~\"[a-z]+\"}",
		expected: []Token{
			{Kind: TokenOpenParen, Value: "{"},
			{Kind: TokenLiteral, Value: "hello"},
			{Kind: TokenOperator, Value: "=~"},
			{Kind: TokenLiteral, Value: "\"[a-z]+\""},
			{Kind: TokenCloseParen, Value: "}"},
		},
	}, {
		name:  "doesn't match regex",
		input: "{hello!~\"[a-z]+\"}",
		expected: []Token{
			{Kind: TokenOpenParen, Value: "{"},
			{Kind: TokenLiteral, Value: "hello"},
			{Kind: TokenOperator, Value: "!~"},
			{Kind: TokenLiteral, Value: "\"[a-z]+\""},
			{Kind: TokenCloseParen, Value: "}"},
		},
	}, {
		name:  "no parens",
		input: "hello=world",
		expected: []Token{
			{Kind: TokenLiteral, Value: "hello"},
			{Kind: TokenOperator, Value: "="},
			{Kind: TokenLiteral, Value: "world"},
		},
	}, {
		name:  "invalid operator",
		input: "{hello=:world}",
		expected: []Token{
			{Kind: TokenOpenParen, Value: "{"},
			{Kind: TokenLiteral, Value: "hello"},
			{Kind: TokenOperator, Value: "="},
			{Kind: TokenLiteral, Value: ":"},
			{Kind: TokenLiteral, Value: "world"},
			{Kind: TokenCloseParen, Value: "}"},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, err := Lex(test.input)
			require.NoError(t, err)
			assert.Equal(t, test.expected, tokens)
		})
	}
}
