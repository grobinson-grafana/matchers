package matchers

import (
	"testing"

	"github.com/prometheus/alertmanager/pkg/labels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected labels.Matchers
		error    string
	}{{
		name:     "no braces",
		input:    "",
		expected: nil,
	}, {
		name:     "open and closing braces",
		input:    "{}",
		expected: nil,
	}, {
		name:     "equals",
		input:    "{foo=\"bar\"}",
		expected: labels.Matchers{mustNewMatcher(t, labels.MatchEqual, "foo", "bar")},
	}, {
		name:     "equals unicode emoji",
		input:    "{foo=\"🙂\"}",
		expected: labels.Matchers{mustNewMatcher(t, labels.MatchEqual, "foo", "🙂")},
	}, {
		name:     "equals without quotes",
		input:    "{foo=bar}",
		expected: labels.Matchers{mustNewMatcher(t, labels.MatchEqual, "foo", "bar")},
	}, {
		name:     "equals without braces",
		input:    "foo=\"bar\"",
		expected: labels.Matchers{mustNewMatcher(t, labels.MatchEqual, "foo", "bar")},
	}, {
		name:     "equals without braces or quotes",
		input:    "foo=bar",
		expected: labels.Matchers{mustNewMatcher(t, labels.MatchEqual, "foo", "bar")},
	}, {
		name:     "equals with trailing comma",
		input:    "{foo=\"bar\",}",
		expected: labels.Matchers{mustNewMatcher(t, labels.MatchEqual, "foo", "bar")},
	}, {
		name:     "equals without braces but trailing comma",
		input:    "foo=\"bar\",",
		expected: labels.Matchers{mustNewMatcher(t, labels.MatchEqual, "foo", "bar")},
	}, {
		name:     "equals with newline",
		input:    "{foo=\"bar\\n\"}",
		expected: labels.Matchers{mustNewMatcher(t, labels.MatchEqual, "foo", "bar\n")},
	}, {
		name:     "equals with tab",
		input:    "{foo=\"bar\\t\"}",
		expected: labels.Matchers{mustNewMatcher(t, labels.MatchEqual, "foo", "bar\t")},
	}, {
		name:     "equals with escaped quotes",
		input:    "{foo=\"\\\"bar\\\"\"}",
		expected: labels.Matchers{mustNewMatcher(t, labels.MatchEqual, "foo", "\"bar\"")},
	}, {
		name:     "equals with escaped backslash",
		input:    "{foo=\"bar\\\\\"}",
		expected: labels.Matchers{mustNewMatcher(t, labels.MatchEqual, "foo", "bar\\")},
	}, {
		name:     "not equals",
		input:    "{foo!=\"bar\"}",
		expected: labels.Matchers{mustNewMatcher(t, labels.MatchNotEqual, "foo", "bar")},
	}, {
		name:     "match regex",
		input:    "{foo=~\"[a-z]+\"}",
		expected: labels.Matchers{mustNewMatcher(t, labels.MatchRegexp, "foo", "[a-z]+")},
	}, {
		name:     "doesn't match regex",
		input:    "{foo!~\"[a-z]+\"}",
		expected: labels.Matchers{mustNewMatcher(t, labels.MatchNotRegexp, "foo", "[a-z]+")},
	}, {
		name:  "complex",
		input: "{foo=\"bar\",bar!=\"baz\"}",
		expected: labels.Matchers{
			mustNewMatcher(t, labels.MatchEqual, "foo", "bar"),
			mustNewMatcher(t, labels.MatchNotEqual, "bar", "baz"),
		},
	}, {
		name:  "complex without quotes",
		input: "{foo=bar,bar!=baz}",
		expected: labels.Matchers{
			mustNewMatcher(t, labels.MatchEqual, "foo", "bar"),
			mustNewMatcher(t, labels.MatchNotEqual, "bar", "baz"),
		},
	}, {
		name:  "complex without braces",
		input: "foo=\"bar\",bar!=\"baz\"",
		expected: labels.Matchers{
			mustNewMatcher(t, labels.MatchEqual, "foo", "bar"),
			mustNewMatcher(t, labels.MatchNotEqual, "bar", "baz"),
		},
	}, {
		name:  "complex without braces or quotes",
		input: "foo=bar,bar!=baz",
		expected: labels.Matchers{
			mustNewMatcher(t, labels.MatchEqual, "foo", "bar"),
			mustNewMatcher(t, labels.MatchNotEqual, "bar", "baz"),
		},
	}, {
		name:  "open brace",
		input: "{",
		error: "0:1: end of input: expected close brace",
	}, {
		name:  "close brace",
		input: "}",
		error: "0:1: }: expected opening brace",
	}, {
		name:  "no open brace",
		input: "foo=\"bar\"}",
		error: "0:10: }: expected opening brace",
	}, {
		name:  "no close brace",
		input: "{foo=\"bar\"",
		error: "0:10: end of input: expected close brace",
	}, {
		name:  "invalid operator",
		input: "{foo=:\"bar\"}",
		error: "5:6: :: invalid input: expected label value",
	}, {
		name:  "another invalid operator",
		input: "{foo%=\"bar\"}",
		error: "4:5: %: invalid input: expected an operator such as '=', '!=', '=~' or '!~'",
	}, {
		name:  "invalid escape sequence",
		input: "{foo=\"bar\\w\"}",
		error: "5:12: \"bar\\w\": invalid input",
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matchers, err := Parse(test.input)
			if test.error != "" {
				require.EqualError(t, err, test.error)
			} else {
				require.Nil(t, err)
				assert.EqualValues(t, test.expected, matchers)
			}
		})
	}
}

func mustNewMatcher(t *testing.T, op labels.MatchType, name, value string) *labels.Matcher {
	m, err := labels.NewMatcher(op, name, value)
	require.NoError(t, err)
	return m
}
