package types

import (
	"testing"

	"github.com/cosiner/golib/test"
)

func TestStringIn(t *testing.T) {
	test.AssertEq(t, "StringIn", 1, StringIn("b", []string{"a", "b", "c"}))
	test.AssertEq(t, "StringIn1", 0, StringIn("a", []string{"a", "b", "c"}))
	test.AssertEq(t, "StringIn2", -1, StringIn("d", []string{"a", "b", "c"}))
}
