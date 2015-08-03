package counter

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestCounter(t *testing.T) {
	tt := testing2.Wrap(t)
	c := New()

	tt.Eq(1, c.Add("a"))
	tt.Eq(2, c.Add("a"))
	tt.Eq(3, c.Add("a"))
	tt.Eq(4, c.Add("a"))

	tt.DeepEq([]string{"a"}, c.Keys())
	tt.Eq(4, c.Count("a"))
	tt.Eq(0, c.Clear("b"))

	tt.Eq(3, c.Remove("a"))
	tt.Eq(2, c.Remove("a"))
	tt.Eq(1, c.Remove("a"))
	tt.Eq(0, c.Remove("a"))
	tt.Eq(0, c.Remove("a"))

}
