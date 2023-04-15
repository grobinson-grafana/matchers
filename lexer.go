package matchers

import (
	"fmt"
	"strings"
	"text/scanner"
)

type TokenKind int

const (
	TokenNone TokenKind = iota
	TokenCloseParen
	TokenComma
	TokenLiteral
	TokenOpenParen
	TokenOperator
)

func (k TokenKind) String() string {
	switch k {
	case TokenCloseParen:
		return "CloseParen"
	case TokenComma:
		return "Comma"
	case TokenLiteral:
		return "Literal"
	case TokenOpenParen:
		return "OpenParen"
	case TokenOperator:
		return "Op"
	default:
		return "None"
	}
}

type Token struct {
	Kind  TokenKind
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("(%s) '%s'", t.Kind, t.Value)
}

func Lex(s string) ([]Token, error) {
	var (
		tokens []Token
		sr     scanner.Scanner
	)
	sr.Init(strings.NewReader(s))

	for tok := sr.Scan(); tok != scanner.EOF; tok = sr.Scan() {
		switch tok {
		case '{':
			tokens = append(tokens, Token{
				Kind:  TokenOpenParen,
				Value: string(tok),
			})
		case '}':
			tokens = append(tokens, Token{
				Kind:  TokenCloseParen,
				Value: string(tok),
			})
		case '=':
			v := string(tok)
			if next := sr.Peek(); next == '~' {
				v = v + string(sr.Next())
			}
			tokens = append(tokens, Token{
				Kind:  TokenOperator,
				Value: v,
			})
		case '!':
			v := string(tok)
			if next := sr.Peek(); next == '=' || next == '~' {
				v = v + string(sr.Next())
			}
			tokens = append(tokens, Token{
				Kind:  TokenOperator,
				Value: v,
			})
		case ',':
			tokens = append(tokens, Token{
				Kind:  TokenComma,
				Value: ",",
			})
		default:
			tokens = append(tokens, Token{
				Kind:  TokenLiteral,
				Value: sr.TokenText(),
			})
		}
	}
	return tokens, nil
}
