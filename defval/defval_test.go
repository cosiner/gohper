package defval

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestInt(t *testing.T) {
	tt := testing2.Wrap(t)
	val := 0
	Int(&val, 10)
	tt.Eq(val, 10)
}

func TestString(t *testing.T) {
	tt := testing2.Wrap(t)
	val := ""
	String(&val, "abc")
	tt.Eq(val, "abc")
}

func TestNil(t *testing.T) {
	var f func()
	var v bool
	Nil(&f, func() {
		v = true
	})
	f()
	testing2.True(t, true)
}

func TestCond(t *testing.T) {
	tt := testing2.Wrap(t)

	tt.Eq(Cond(true).String("a", "b"), "a")
}
