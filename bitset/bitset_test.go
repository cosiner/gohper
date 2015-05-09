package bitset

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestZeroLength(t *testing.T) {
	tt := testing2.Wrap(t)
	defer tt.Recover()
	NewBitset(0)
}

func TestBitSet(t *testing.T) {
	tt := testing2.Wrap(t)
	list := []uint{1, 2, 4, 5, 6, 7, 60}
	s := NewBitset(32, 1)
	tt.True(s.UnitCount() == 1)
	tt.True(s.Uint() == 2)
	tt.True(s.Uint64() == 2)
	testBitset(tt, s, list)
}

func testBitset(tt testing2.Test, s *Bitset, list []uint) {
	s.UnsetAll()

	for _, l := range list {
		s.Set(l)
		tt.True(s.IsSet(l))

		s.Unset(l)
		tt.False(s.IsSet(l))

		s.SetTo(l, true)
		tt.True(s.IsSet(l))

		s.SetTo(l, false)
		tt.False(s.IsSet(l))

		s.Flip(l)
		tt.True(s.IsSet(l))

		s.FlipAll()
		tt.False(s.IsSet(l))

		s.SetAll()
		tt.True(s.IsSet(l))

		s.Unset(l)
		tt.Eq(s.BitCount(), int(s.Length(0)-1))

		s.UnsetAll()
		tt.False(s.IsSet(l))
	}
	tt.Eq(uint(64), s.UnitLen())
	tt.Eq(uint(1), s.UnitCount())

	s.UnsetAll()
	for _, l := range list {
		s.Set(l)
	}
	tt.DeepEq(s.Bits(), list)
	s.Flip(127)
	tt.Eq(uint(128), s.Length(0))

	tt.True(s.Length(64) == 64)
	s.UnsetAll()

	s.Except(list...)
	for _, l := range list {
		tt.False(s.IsSet(l))
	}

	cl := s.Clone()
	cl.Length(256)
	tt.True(cl.UnitCount() == 4)
	cl.Except(list...).Except(cl.Bits()...)
	s.Intersection(cl)
	tt.True(s.BitCount() == 0)

	s.Union(cl)
	tt.True(s.BitCount() == len(list))

	s.Diff(cl)
	tt.True(s.BitCount() == 0)

	s.UnsetAll()
	cl.UnsetAll()
	cl.Intersection(s)
	tt.True(cl.BitCount() == 0)
}
