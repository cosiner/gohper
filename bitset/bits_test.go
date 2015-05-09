package bitset

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestBits(t *testing.T) {
	tt := testing2.Wrap(t)
	list := []uint{2, 4, 5, 6, 7, 60}
	s1 := NewBits()
	s2 := BitsList(list...)

	var u uint
	for _, l := range list {
		u |= 1 << l
	}
	s3 := BitsFrom(u)

	testBits(tt, s1, list, true)
	testBits(tt, s2, list, false)
	testBits(tt, s3, list, false)
}

func testBits(tt testing2.Test, s *Bits, list []uint, empty bool) {
	for _, l := range list {
		s.Set(l)
		tt.True(s.IsSet(l))

		s.Unset(l)
		tt.False(s.IsSet(l))

		s.SetTo(l, true)
		tt.True(s.IsSet(l))

		if empty {
			tt.True(s.Uint() == 1<<l)
			tt.True(s.Uint64() == 1<<l)
		}

		s.SetTo(l, false)
		tt.False(s.IsSet(l))

		s.Flip(l)
		tt.True(s.IsSet(l))

		s.FlipAll()
		tt.False(s.IsSet(l))

		s.SetAll()
		tt.True(s.IsSet(l))

		s.UnsetAll()
		tt.False(s.IsSet(l))

		s.SetBefore(l)
		for i := uint(0); i < l; i++ {
			tt.True(s.IsSet(i))
		}

		s.SetSince(l)
		for i := l; i < 64; i++ {
			tt.True(s.IsSet(uint(i)))
		}

		s.Unset(l)
		tt.True(s.BitCount() == 63)

		s.UnsetBefore(l)
		for i := uint(0); i < l; i++ {
			tt.False(s.IsSet(i))
		}

		s.UnsetSince(l)
		for i := l; i < 64; i++ {
			tt.False(s.IsSet(i))
		}
	}
}

func TestBitCount(t *testing.T) {
	tt := testing2.Wrap(t)
	list := []uint{2, 4, 5, 6, 7, 60}
	var u uint
	for _, l := range list {
		u |= 1 << l
	}

	tt.True(BitCountUint(u) == len(list))
	tt.True(BitCount(uint64(u)) == len(list))
}
