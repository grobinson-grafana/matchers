package matchers

type Iterator struct {
	lexer   Lexer
	current Token
	next    Token
}

func NewIterator(lexer Lexer) Iterator {
	return Iterator{
		lexer: lexer,
	}
}

func (it *Iterator) Next() Token {
	if it.next != (Token{}) {
		it.current = it.next
		it.next = Token{}
		return it.current
	}
	if tok, err := it.lexer.Scan(); err == nil && tok.Kind != TokenNone {
		it.current = tok
		return it.current
	}
	return Token{Kind: TokenNone}
}

func (it *Iterator) Peek() Token {
	if it.next == (Token{}) {
		if tok, err := it.lexer.Scan(); err == nil && tok.Kind != TokenNone {
			it.next = tok
		} else {
			it.next = Token{Kind: TokenNone}
		}
	}
	return it.next
}
