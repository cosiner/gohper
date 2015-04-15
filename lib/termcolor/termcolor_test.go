package termcolor

import (
	"testing"
)

func TestColor(t *testing.T) {
	tc := NewColor().Bg("green").Highlight().Inverse().Underline().Fg(RED)
	t.Log(tc.Render("aaa"))
	t.Log(tc.Render("aaadd"))
}
