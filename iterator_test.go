package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	helloExample       = "{hello=\"world\"}"
	helloExampleTokens = []Token{{
		Kind:  TokenOpenParen,
		Value: "{",
	}, {
		Kind:  TokenIdent,
		Value: "hello",
	}, {
		Kind:  TokenOperator,
		Value: "=",
	}, {
		Kind:  TokenQuoted,
		Value: "\"world\"",
	}, {
		Kind:  TokenCloseParen,
		Value: "}",
	}}
)

func TestIterator_Next(t *testing.T) {
	it := NewIterator(NewLexer(helloExample))
	for _, expected := range helloExampleTokens {
		assert.Equal(t, expected, it.Next())
	}
	// check that end of iterator returns Token{}
	assert.Equal(t, Token{}, it.Next())
	// check that iterator returns Token{} past the end
	assert.Equal(t, Token{}, it.Next())
}

func TestIterator_Peek(t *testing.T) {
	it := NewIterator(NewLexer(helloExample))
	expected := helloExampleTokens
	// check can peek first item
	assert.Equal(t, expected[0], it.Peek())
	// check successive peeks don't advance the iterator
	assert.Equal(t, expected[0], it.Peek())
	// check that advance iterator returns peeked item
	assert.Equal(t, expected[0], it.Next())
	// check that peek returns second item
	assert.Equal(t, expected[1], it.Peek())
	// check the peek and next work together
	for i := 1; i < len(expected); i++ {
		assert.Equal(t, expected[i], it.Peek())
		assert.Equal(t, expected[i], it.Next())
	}
	// check that end of iterator returns Token{}
	assert.Equal(t, Token{}, it.Peek())
	// check that iterator returns Token{} past end
	assert.Equal(t, Token{}, it.Next())
	assert.Equal(t, Token{}, it.Peek())
}
