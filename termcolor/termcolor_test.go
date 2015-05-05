package termcolor

import (
	"os"
	"testing"
)

func TestRender(t *testing.T) {
	tc := Blue.Finish()
	t.Log(tc.Render("aaa"))
	t.Log(tc.Render("aaadd"))
}

func TestRenderTo(t *testing.T) {
	Red.RenderTo(os.Stdout, "aaaaaaaaaaaaa\n")
	Red.RenderTo(os.Stdout, "bbbbb\n")
}
