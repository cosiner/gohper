package sync2

import (
	"testing"
	"github.com/cosiner/gohper/testing2"
)

func TestFlag(t *testing.T) {
	tt := testing2.Wrap(t)

	var flag Flag
	tt.False(flag.IsTrue())
	tt.True(flag.MakeTrue())
	tt.False(flag.MakeTrue())
	tt.True(flag.MakeFalse())
	tt.False(flag.MakeFalse())

	var flags Flags
	tt.False(flags.IsTrue("a"))
	tt.True(flags.MakeTrue("a"))
	tt.False(flags.MakeTrue("a"))
	tt.True(flags.MakeFalse("a"))
	tt.False(flags.MakeFalse("a"))
}