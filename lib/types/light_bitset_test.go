package types

import (
	"github.com/cosiner/gohper/lib/test"

	"testing"
)

func TestLightBitSet(t *testing.T) {
	tt := test.Wrap(t)
	l := NewLightBitSet()
	l.Set(2)
	l.Set(4)
	l.Set(60)
	tt.True(l.IsSet(2))
	tt.True(l.IsSet(4))
	tt.True(l.IsSet(60))
	l.Unset(60)
	tt.False(l.IsSet(60))
	l.Flip(60)
	tt.True(l.IsSet(60))
	l.SetTo(60, false)
	tt.False(l.IsSet(60))
	l.FlipAll()
	tt.True(l.IsSet(61))
}

func TestSinceBefore(t *testing.T) {
	tt := test.Wrap(t)
	l := NewLightBitSet()
	l.SetAllBefore(9)
	for i := 0; i < 9; i++ {
		tt.True(l.IsSet(uint(i)))
	}
	l.UnsetAll()
	l.SetAllSince(9)
	for i := 9; i < 64; i++ {
		tt.True(l.IsSet(uint(i)))
	}
	l.SetAll()
	l.UnsetAllBefore(9)
	for i := 0; i < 9; i++ {
		tt.False(l.IsSet(uint(i)))
	}
	l.SetAll()
	l.UnsetAllSince(9)
	for i := 0; i < 9; i++ {
		tt.True(l.IsSet(uint(i)))
	}
	for i := 9; i < 64; i++ {
		tt.False(l.IsSet(uint(i)))
	}
}

func BenchmarkBitCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = BitCountUint(0x121241112a)
	}
}

func TestBitCount(t *testing.T) {
	tt := test.Wrap(t)
	tt.Eq(2, BitCountUint(3))
	tt.Log(BitCount(3))
}
