package matchers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/prometheus/alertmanager/pkg/labels"
)

var (
	ErrEOF            = errors.New("EOF")
	ErrNoOpeningParen = errors.New("expected opening '{'")
	ErrNoClosingParen = errors.New("expected closing '}'")
)

func Parse(input string) (labels.Matchers, error) {
	var (
		tok      Token
		matchers labels.Matchers
	)

	l := NewLexer(input)

	// Must start with opening paren
	if tok, _ = l.Scan(); tok.Kind != TokenOpenParen {
		return nil, ErrNoOpeningParen
	}

	for {
		// Break if there is a closing paren
		if tok, _ = l.Peek(); tok.Kind == TokenCloseParen {
			break
		}

		// The next token should be the label name
		if tok, _ = l.Scan(); tok.Kind != TokenIdent {
			return nil, fmt.Errorf("expected a label name, got '%s'", tok.Value)
		}
		name := tok.Value

		// The next token after the label name should be one of the expected operators
		if tok, _ = l.Scan(); tok.Kind != TokenOperator {
			return nil, fmt.Errorf("expected one of '=', '!=', '=~' or '!~', got '%s'", tok.Value)
		}
		op := labels.MatchEqual
		switch tok.Value {
		case "=":
			op = labels.MatchEqual
		case "!=":
			op = labels.MatchNotEqual
		case "=~":
			op = labels.MatchRegexp
		case "!~":
			op = labels.MatchNotRegexp
		default:
			return nil, fmt.Errorf("expected one of '=', '!=', '=~' or '!~', got '%s'", tok.Value)
		}

		// The next token after the operator should be the label value
		if tok, _ = l.Scan(); tok.Kind != TokenQuoted {
			return nil, fmt.Errorf("expected a label value, got '%s'", tok.Value)
		}
		value := strings.TrimPrefix(strings.TrimSuffix(tok.Value, "\""), "\"")

		m, err := labels.NewMatcher(op, name, value)
		if err != nil {
			return nil, fmt.Errorf("failed to create matcher: %s", err)
		}
		matchers = append(matchers, m)

		// If the next token is not a comma then it has to be a closing paren
		if tok, _ = l.Peek(); tok.Kind != TokenComma && tok.Kind != TokenCloseParen {
			return nil, fmt.Errorf("expected comma or closing '}', got '%s'", tok.Value)
		}

		// If the next token is a comma then expect more matchers
		if tok, _ = l.Peek(); tok.Kind == TokenComma {
			l.Scan()
			// The next token is a comma so the one after that must be a label name
			if tok, _ = l.Peek(); tok.Kind != TokenIdent {
				return nil, fmt.Errorf("expected label name after comma, got '%s'", tok.Value)
			}
		}
	}

	// Must end with closing paren
	if tok, _ = l.Scan(); tok.Kind != TokenCloseParen {
		return nil, ErrNoClosingParen
	}

	// There should be no more tokens
	if tok, _ = l.Scan(); tok != (Token{}) {
		return nil, fmt.Errorf("expected end of input, got '%s'", tok.Value)
	}

	return matchers, nil
}
