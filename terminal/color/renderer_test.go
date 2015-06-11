package color

import (
	"os"
	"testing"
)

func TestRender(t *testing.T) {
	r := LightBlue
	t.Log(r.Render("aaa"))
	t.Log(r.Render("aaadd"))
}

func TestRenderTo(t *testing.T) {
	LightRed.RenderTo(os.Stdout, "aaaaaaaaaaaaa\n")
	LightRed.RenderTo(os.Stdout, "bbbbb\n")
}
