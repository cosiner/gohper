package color

import (
	"testing"
)

func TestRender(t *testing.T) {
	r := New(Highlight, FgGreen)
	t.Log(r.RenderString("aaa"))
	t.Log(r.RenderString("aaadd"))
}
