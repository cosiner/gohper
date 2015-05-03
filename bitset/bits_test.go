package bitset

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestLightBitSet(t *testing.T) {
	tt := testing2.Wrap(t)
	s := NewBits()
	s.Set(2)
	s.Set(4)
	s.Set(60)
	tt.True(s.IsSet(2))
	tt.True(s.IsSet(4))
	tt.True(s.IsSet(60))
	s.Unset(60)
	tt.False(s.IsSet(60))
	s.Flip(60)
	tt.True(s.IsSet(60))
	s.SetTo(60, false)
	tt.False(s.IsSet(60))
	s.FlipAll()
	tt.True(s.IsSet(61))
}

func TestSinceBefore(t *testing.T) {
	tt := testing2.Wrap(t)
	s := NewBits()
	s.SetBefore(9)
	for i := 0; i < 9; i++ {
		tt.True(s.IsSet(uint(i)))
	}
	s.UnsetAll()
	s.SetSince(9)
	for i := 9; i < 64; i++ {
		tt.True(s.IsSet(uint(i)))
	}
	s.SetAll()
	s.UnsetBefore(9)
	for i := 0; i < 9; i++ {
		tt.False(s.IsSet(uint(i)))
	}
	s.SetAll()
	s.UnsetSince(9)
	for i := 0; i < 9; i++ {
		tt.True(s.IsSet(uint(i)))
	}
	for i := 9; i < 64; i++ {
		tt.False(s.IsSet(uint(i)))
	}
}

func BenchmarkBitCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = BitCountUint(0x121241112a)
	}
}

func TestBitCount(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.Eq(2, BitCountUint(3))
	tt.Log(BitCount(3))
}
