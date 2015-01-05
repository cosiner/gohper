package test

import (
	"testing"
)

func TestTesttool(t *testing.T) {
	AssertEq(t, true, false, "AssertEq")
	AssertEq(t, true, true, "AssertEq")
	AssertEq(t, true, 1 <= 2, "AssertEq")
}
