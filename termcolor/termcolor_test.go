package termcolor

import (
	"os"
	"testing"
)

func TestRender(t *testing.T) {
	tc := New().Bg(GREEN).Highlight().Inverse().Underline().Fg(RED).Finish()
	t.Log(tc.Render("aaa"))
	t.Log(tc.Render("aaadd"))
}

func TestRenderTo(t *testing.T) {
	tc := New().Bg(BLUE).Finish()
	tc.RenderTo(os.Stdout, "aaaaaaaaaaaaa\n")
	tc.RenderTo(os.Stdout, "bbbbb\n")
}
