package matchers

import (
	"testing"
)

const (
	simpleExample  = "{foo=bar}"
	complexExample = "{foo=bar,bar=~\"[a-zA-Z0-9+]\"}"
)

func BenchmarkLexSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := Lex(simpleExample); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLexComplex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := Lex(complexExample); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := Parse(simpleExample); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseComplex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := Parse(complexExample); err != nil {
			b.Fatal(err)
		}
	}
}
