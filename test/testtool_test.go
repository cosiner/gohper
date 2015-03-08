package test

import (
	"testing"
)

func TestTesttool(t *testing.T) {
	AssertTrue(t, false)
	AssertEq(t, true, false)
	AssertNE(t, true, 1 > 2)
}
