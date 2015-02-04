package types

import (
	"github.com/cosiner/golib/test"

	"testing"
)

func BenchmarkLightBitSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := NewLightBitSet()
		l.SetFor(2, 4, 6, 10, 13)
		// _ = bs.IsSet(4)
	}
}

func TestLightBitSet(t *testing.T) {
	tt := test.WrapTest(t)
	l := NewLightBitSet()
	l.Set(2)
	l.Set(4)
	l.Set(60)
	tt.AssertTrue("1", l.IsSet(2))
	tt.AssertTrue("1", l.IsSet(4))
	tt.AssertTrue("1", l.IsSet(60))
	l.Unset(60)
	tt.AssertFalse("1", l.IsSet(60))
	l.Flip(60)
	tt.AssertTrue("1", l.IsSet(60))
	l.SetTo(60, false)
	tt.AssertFalse("1", l.IsSet(60))
	l.FlipAll()
	tt.AssertTrue("1", l.IsSet(61))
}
