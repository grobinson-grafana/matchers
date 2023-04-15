package matchers

type Iterator struct {
	tokens []Token
	pos    int
}

func NewIterator(tokens []Token) Iterator {
	return Iterator{tokens: tokens, pos: -1}
}

func (it *Iterator) IsLast() bool {
	return it.pos == len(it.tokens)-1
}

func (it *Iterator) Next() Token {
	if !it.IsLast() {
		it.pos += 1
		return it.tokens[it.pos]
	}
	return Token{}
}

func (it *Iterator) Peek() Token {
	if !it.IsLast() {
		return it.tokens[it.pos+1]
	}
	return Token{}
}

func (it *Iterator) Pos() int {
	return it.pos
}
