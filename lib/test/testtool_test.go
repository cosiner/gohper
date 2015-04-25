package test

import "testing"

func TestTesttool(t *testing.T) {
	Eq(t, 1, 1)
	NE(t, t, nil)
}
