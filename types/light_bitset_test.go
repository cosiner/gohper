package types

import (
	"github.com/cosiner/golib/test"

	"testing"
)

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

func TestSinceBefore(t *testing.T) {
	tt := test.WrapTest(t)
	l := NewLightBitSet()
	l.SetAllBefore(9)
	for i := 0; i < 9; i++ {
		tt.AssertTrue("1", l.IsSet(uint(i)))
	}
	l.UnsetAll()
	l.SetAllSince(9)
	for i := 9; i < 64; i++ {
		tt.AssertTrue("1", l.IsSet(uint(i)))
	}
	l.SetAll()
	l.UnsetAllBefore(9)
	for i := 0; i < 9; i++ {
		tt.AssertFalse("1", l.IsSet(uint(i)))
	}
	l.SetAll()
	l.UnsetAllSince(9)
	for i := 0; i < 9; i++ {
		tt.AssertTrue("1", l.IsSet(uint(i)))
	}
	for i := 9; i < 64; i++ {
		tt.AssertFalse("1", l.IsSet(uint(i)))
	}
}

func BenchmarkBitCount(b *testing.B) {
	l := NewLightBitSet()
	l.SetAll()
	for i := 0; i < b.N; i++ {
		_ = l.BitCount()
	}
}
