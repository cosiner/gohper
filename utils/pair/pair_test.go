package pair

import (
	"strings"
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

	p.Trim()
	tt.Eq("(aa:false)", p.String())

	p = Parse("'aa'='false'", "=")
	p.TrimQuote()
	tt.Eq("(aa:false)", p.String())
}

func TestRparse(t *testing.T) {
	tt := testing2.Wrap(t)
	p := Rparse("aa=bb=1", "=")

	val, err := p.IntValue()
	tt.Nil(err)
	tt.Eq(1, val)
}

func TestPair2(t *testing.T) {
	tt := testing2.Wrap(t)
	p := Parse("=", "=")

	tt.True(p.NoKey())
	tt.True(p.NoValue())

	p = Parse("abc", "=")
	tt.True(p.NoValue())
	tt.Eq(p.Key, "abc")
}

func TestPair3(t *testing.T) {
	tt := testing2.Wrap(t)

	p := ParsePairWith("a=b", "=", strings.Index)
	tt.True(p.HasKey())
	tt.True(p.HasValue())
	tt.Eq("b", p.ValueOrKey())

	p = Parse("a=", "=")
	tt.Eq("a", p.ValueOrKey())
	p = Parse("a", "=")
	tt.Eq("a", p.ValueOrKey())

	p = Parse("'a=b", "=")
	tt.False(p.TrimQuote())
	p = Parse("a=b'", "=")
	tt.Log(p.Value)
	tt.False(p.TrimQuote())
}
