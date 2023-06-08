package matchers

import (
	"fmt"
)

type TokenKind int

const (
	TokenNone TokenKind = iota
	TokenCloseBrace
	TokenComma
	TokenIdent
	TokenOpenBrace
	TokenOperator
	TokenQuoted
)

func (k TokenKind) String() string {
	switch k {
	case TokenCloseBrace:
		return "CloseBrace"
	case TokenComma:
		return "Comma"
	case TokenIdent:
		return "Ident"
	case TokenOpenBrace:
		return "OpenBrace"
	case TokenOperator:
		return "Op"
	case TokenQuoted:
		return "Quoted"
	default:
		return "None"
	}
}

type Token struct {
	Kind  TokenKind
	Value string
	Position
}

func (t Token) String() string {
	return fmt.Sprintf("(%s) '%s'", t.Kind, t.Value)
}

func IsNone(t Token) bool {
	return t == Token{}
}

type Position struct {
	OffsetStart int // The start position in the input
	OffsetEnd   int // The end position in the input
	ColumnStart int // The column number
	ColumnEnd   int // The end of the column
}
