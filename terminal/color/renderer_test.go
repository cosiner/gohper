package color

import (
	"testing"

	"github.com/cosiner/gohper/bytes2"
)

func TestRender(t *testing.T) {
	r := LightBlue
	t.Log(r.Render("aaa"))
	t.Log(r.Render("aaadd"))
}

func TestRenderTo(t *testing.T) {
	buf := bytes2.MakeBuffer(0, 128)
	LightRed.RenderTo(buf, "aaaaaaaaaaaaa\n")
	LightRed.RenderTo(buf, "bbbbb\n")
}
