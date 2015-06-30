package unibyte

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestLetter(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.True(IsLower('z'))
	tt.True(IsUpper('Z'))
	tt.True(IsLetter('a'))
	tt.True(IsSpaceQuote('\''))

	tt.True('a' == ToLower('A'))
	tt.True('A' == ToUpper('a'))

	tt.Eq("a", ToLowerString('A'))
	tt.Eq("A", ToUpperString('a'))
}
