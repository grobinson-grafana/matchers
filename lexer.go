package matchers

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	eof rune = -1
)

// Lexer scans a sequence of tokens that match the grammar of Prometheus-like
// matchers. A token is emitted for each call to Scan() which returns the
// next token in the input or an error if the input does not conform to the
// grammar. A token can be one of a number of kinds and corresponds to a
// subslice of the input. Once the input has been consumed successive calls to
// Scan() return a TokenNone token.
type Lexer struct {
	input string
	err   error
	start int // the start of the current token
	pos   int // the position of the cursor in the input
	width int // the width of the last rune
}

func NewLexer(input string) Lexer {
	return Lexer{
		input: input,
	}
}

func (l *Lexer) Scan() (Token, error) {
	// Do not attempt to emit more tokens if the input is invalid
	if l.err != nil {
		return Token{}, l.err
	}

	// Iterate over each rune in the input and either emit a token or an error
	for r := l.next(); r != eof; r = l.next() {
		switch {
		case r == '{':
			return l.emit(TokenOpenParen), nil
		case r == '}':
			return l.emit(TokenCloseParen), nil
		case r == ',':
			return l.emit(TokenComma), nil
		case r == '=' || r == '!':
			l.rewind()
			var tok Token
			tok, l.err = l.scanOperator()
			return tok, l.err
		case r == '"':
			l.rewind()
			var tok Token
			tok, l.err = l.scanQuoted()
			return tok, l.err
		case unicode.IsLetter(r):
			l.rewind()
			var tok Token
			tok, l.err = l.scanIdent()
			return tok, l.err
		case unicode.IsSpace(r):
			l.skip()
		default:
			l.err = fmt.Errorf("unexpected input: %s", l.input[0:l.pos])
			return Token{}, l.err
		}
	}

	return Token{}, nil
}

func (l *Lexer) scanIdent() (Token, error) {
	for r := l.next(); r != eof; r = l.next() {
		if !unicode.IsLetter(r) {
			l.rewind()
			break
		}
	}
	return l.emit(TokenIdent), nil
}

func (l *Lexer) scanOperator() (Token, error) {
	l.acceptRun("!=~")
	return l.emit(TokenOperator), nil
}

func (l *Lexer) scanQuoted() (Token, error) {
	if err := l.expect("\""); err != nil {
		return Token{}, err
	}

	var isEscaped bool
	for r := l.next(); r != eof; r = l.next() {
		if r == '\\' {
			isEscaped = true
		} else if r == '"' && !isEscaped {
			l.rewind()
			break
		} else {
			isEscaped = false
		}
	}

	if err := l.expect("\""); err != nil {
		return Token{}, err
	}

	return l.emit(TokenQuoted), nil
}

func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.rewind()
	return false
}

func (l *Lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.rewind()
}

func (l *Lexer) expect(valid string) error {
	r := l.next()
	if r == -1 {
		return fmt.Errorf("expected one of '%s', got EOF", valid)
	}
	if strings.IndexRune(valid, r) < 0 {
		return fmt.Errorf("expected one of '%s', got '%c'", valid, r)
	}
	return nil
}

func (l *Lexer) emit(kind TokenKind) Token {
	tok := Token{
		Kind:  kind,
		Value: l.input[l.start:l.pos],
	}
	l.start = l.pos
	return tok
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width
	l.pos += width
	return r
}

func (l *Lexer) rewind() {
	l.pos -= l.width
	l.width = 0
}

func (l *Lexer) skip() {
	l.start = l.pos
}
