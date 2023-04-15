package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	helloWorldExample = []Token{
		{Kind: TokenOpenParen, Value: "{"},
		{Kind: TokenLiteral, Value: "hello"},
		{Kind: TokenOperator, Value: "="},
		{Kind: TokenLiteral, Value: "world"},
		{Kind: TokenCloseParen, Value: "}"},
	}
)

func TestIterator_IsLast(t *testing.T) {
	it := NewIterator(helloWorldExample)
	for i := 0; i < len(helloWorldExample)-1; i++ {
		it.Next()
		assert.False(t, it.IsLast())
	}
	// last token in iterator
	it.Next()
	assert.True(t, it.IsLast())
}

func TestIterator_Next(t *testing.T) {
	it := NewIterator(helloWorldExample)
	assert.Equal(t, it.Pos(), -1)
	for i := 0; i < len(helloWorldExample); i++ {
		// move the iterator to the next token
		assert.Equal(t, helloWorldExample[i], it.Next())
	}
	// end of the iterator
	assert.Equal(t, Token{}, it.Next())
	// pos should not go off the end
	assert.Equal(t, len(helloWorldExample)-1, it.Pos())
}

func TestIterator_Peek(t *testing.T) {
	it := NewIterator(helloWorldExample)
	assert.Equal(t, helloWorldExample[0], it.Next())
	assert.Equal(t, 0, it.Pos())
	// peek the next item
	assert.Equal(t, helloWorldExample[1], it.Peek())
	// pos should not change
	assert.Equal(t, 0, it.Pos())
	// end of the iterator
	for tok := it.Next(); tok != (Token{}); tok = it.Next() {
		// consume all tokens
	}
	assert.Equal(t, Token{}, it.Peek())
}
