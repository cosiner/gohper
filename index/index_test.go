package index

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestStringIn(t *testing.T) {
	testing2.Eq(t, 1, StringIn("b", []string{"a", "b", "c"}))
	testing2.Eq(t, 0, StringIn("a", []string{"a", "b", "c"}))
	testing2.Eq(t, -1, StringIn("d", []string{"a", "b", "c"}))
}

func TestCharIn(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.Log(CharIn('a', "abc"))
	tt.Log(CharIn('a', "abcd"))

	tt.Log(CharIn('b', "abc"))
	tt.Log(CharIn('b', "abcd"))

	tt.Log(CharIn('d', "abd"))
	tt.Log(CharIn('d', "abcd"))

	tt.Log(CharIn('e', "abc"))
	tt.Log(CharIn('e', "abcd"))
}
