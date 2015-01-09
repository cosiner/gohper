package test

import (
	"testing"
)

func TestTesttool(t *testing.T) {
	AssertTrue(t, "AssertEq", true)
	AssertEq(t, "AssertEq", true, true)
	AssertNE(t, "AssertEq", true, 1 == 2)
}
