package types

import (
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

func TestStringIn(t *testing.T) {
	test.AssertEq(t, 1, StringIn("b", []string{"a", "b", "c"}))
	test.AssertEq(t, 0, StringIn("a", []string{"a", "b", "c"}))
	test.AssertEq(t, -1, StringIn("d", []string{"a", "b", "c"}))
}

func TestCharIn(t *testing.T) {
	tt := test.WrapTest(t)
	tt.Log(CharIn('a', "abc"))
	tt.Log(CharIn('a', "abcd"))

	tt.Log(CharIn('b', "abc"))
	tt.Log(CharIn('b', "abcd"))

	tt.Log(CharIn('d', "abd"))
	tt.Log(CharIn('d', "abcd"))

	tt.Log(CharIn('e', "abc"))
	tt.Log(CharIn('e', "abcd"))
}
