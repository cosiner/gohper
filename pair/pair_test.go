package pair

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestPair(t *testing.T) {
	tt := testing2.Wrap(t)
	p := Parse("aa=false", "=")
	tt.True(p.HasKey())
	v, e := p.BoolValue()
	tt.True(e == nil)
	tt.False(v)
}
