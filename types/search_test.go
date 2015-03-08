package types

import (
	"testing"

	"github.com/cosiner/golib/test"
)

func TestStringIn(t *testing.T) {
	test.AssertEq(t, 1, StringIn("b", []string{"a", "b", "c"}))
	test.AssertEq(t, 0, StringIn("a", []string{"a", "b", "c"}))
	test.AssertEq(t, -1, StringIn("d", []string{"a", "b", "c"}))
}
