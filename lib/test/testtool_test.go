package test

import "testing"

func TestTesttool(t *testing.T) {
	Nil(t, nil)
	Eq(t, 1, 1)
	NE(t, t, nil)
}
