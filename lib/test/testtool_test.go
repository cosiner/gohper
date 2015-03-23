package test

import "testing"

func TestTesttool(t *testing.T) {
	assertNil(t, 0, nil)
	assertEq(t, 0, 1, 1)
	assertNE(t, 0, t, nil)
}
