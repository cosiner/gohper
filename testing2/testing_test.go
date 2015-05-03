package testing2

import "testing"

func TestTest(t *testing.T) {
	Eq(t, 1, 1)
	NE(t, t, nil)
}
