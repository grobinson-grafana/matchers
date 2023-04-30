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
		name:     "open and closing parens",
		input:    "{}",
		expected: nil,
	}, {
		name:     "equals",
		input:    "{foo=bar}",
		expected: labels.Matchers{mustNewMatcher(t, labels.MatchEqual, "foo", "bar")},
	}, {
		name:     "not equals",
		input:    "{foo!=bar}",
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
		name:  "equals and a not equals",
		input: "{foo=bar,bar!=baz}",
		expected: labels.Matchers{
			mustNewMatcher(t, labels.MatchEqual, "foo", "bar"),
			mustNewMatcher(t, labels.MatchNotEqual, "bar", "baz"),
		},
	}, {
		name:  "equals unicode emoji",
		input: "{foo=\"🙂\"}",
		expected: labels.Matchers{
			mustNewMatcher(t, labels.MatchEqual, "foo", "🙂"),
		},
	}, {
		name:  "no parens",
		input: "foo=bar",
		error: "expected opening '{'",
	}, {
		name:  "invalid operator",
		input: "{foo=:bar}",
		error: "expected comma or closing '}', got 'bar'",
	}, {
		name:  "another invalid operator",
		input: "{foo%=bar}",
		error: "expected one of '=', '!=', '=~' or '!~', got '%'",
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
