package termcolor

import (
	"testing"
)

func TestColor(t *testing.T) {
	tc := NewColor().Bg("green")
	t.Log(tc.Render("aaa"))
}
