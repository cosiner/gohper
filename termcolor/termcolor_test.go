package termcolor

import (
	"os"
	"testing"
)

func TestRender(t *testing.T) {
	tc := LightBlue.Finish()
	t.Log(tc.Render("aaa"))
	t.Log(tc.Render("aaadd"))
}

func TestRenderTo(t *testing.T) {
	LightRed.RenderTo(os.Stdout, "aaaaaaaaaaaaa\n")
	LightRed.RenderTo(os.Stdout, "bbbbb\n")
}
