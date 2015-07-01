// Package bitset implements a Bitset and a small Bits.
package bitset

// u_1 is uint 1
const (
	Uint0       uint   = 0
	unitLenLogN        = 6 // log64 = 6
	unitLen            = 1 << unitLenLogN
	unitMax     uint64 = 1<<unitLen - 1
)

// newUnits create a new unit set has given unit count for bitset
func newUnits(count uint) []uint64 {
	return make([]uint64, count)
}

// unitCount return unit count need for the length
func unitCount(length uint) uint {
	count := length >> unitLenLogN
	if length&(unitLen-1) != 0 {
		count++
	}

	return count
}

// unitPos return the unit position that index bit in
func unitPos(index uint) uint {
	return index >> unitLenLogN
}

// unitIndex return the unit index that index bit in
func unitIndex(index uint) uint {
	return index & (unitLen - 1)
}

// Bitset is a bitset based on a uint64 array
type Bitset struct {
	length uint     // bitset length
	set    []uint64 // bitset data store
}

func nonZeroLen(l uint) {
	if l == 0 {
		panic("length should not be zero")
	}
}

// NewBitset return a new bitset with given length, all index in list are set to 1
func NewBitset(length uint, indexs ...uint) *Bitset {
	nonZeroLen(length)
	s := &Bitset{
		length: length,
		set:    newUnits(unitCount(length)),
	}

	for _, i := range indexs {
		s.Set(i)
	}

	return s
}

// Set set index bit to 1, expand the bitset if index out of length range
func (s *Bitset) Set(index uint) *Bitset {
	s.extend(index)
	s.set[unitPos(index)] |= 1 << unitIndex(index)

	return s
}

// SetAll set all bits to 1
func (s *Bitset) SetAll() *Bitset {
	return s.unitOp(func(s *Bitset, i int) {
		s.set[i] = unitMax
	})
}

// Unset set index bit to 0, expand the bitset if index out of length range
func (s *Bitset) Unset(index uint) *Bitset {
	s.extend(index)
	s.set[unitPos(index)] &= ^(1 << unitIndex(index))

	return s
}

// UnsetAll set all bits to 0
func (s *Bitset) UnsetAll() *Bitset {
	return s.unitOp(func(s *Bitset, i int) {
		s.set[i] = 0
	})
}

// Flip the index bit, expand the bitset if index out of length range
func (s *Bitset) Flip(index uint) *Bitset {
	if index >= s.Length(0) {
		s.Set(index)
	}
	s.set[unitPos(index)] ^= 1 << unitIndex(index)

	return s
}

// FlipAll flip all the index bit
func (s *Bitset) FlipAll() *Bitset {
	return s.unitOp(func(s *Bitset, i int) {
		s.set[i] = ^s.set[i]
	})
}

// Except set all bits except given index to 1, the except bits set to 0
func (s *Bitset) Except(index ...uint) *Bitset {
	s.SetAll()
	for _, i := range index {
		s.Unset(i)
	}

	return s
}

// IsSet check whether or not index bit is set
func (s *Bitset) IsSet(index uint) bool {
	return index < s.Length(0) && (s.set[unitPos(index)]&(1<<unitIndex(index))) != 0
}

// SetTo set index bit to 1 if value is true, otherwise 0
func (s *Bitset) SetTo(index uint, value bool) *Bitset {
	if value {
		return s.Set(index)
	}

	return s.Unset(index)
}

// unitOp iter the bitset unit, apply function to each unit
func (s *Bitset) unitOp(f func(*Bitset, int)) *Bitset {
	for i, n := 0, len(s.set); i < n; i++ {
		f(s, i)
	}

	return s
}

// Union union another bitset to current bitset, expand the bitset if index out of length range
func (s *Bitset) Union(b *Bitset) *Bitset {
	return s.bitsetOp(
		b,

		func(s, b *Bitset, length *uint) {
			bl, l := s.Length(0), *length

			if bl < l {
				s.unsetTop()
				s.Length(l)
			} else if bl > l {
				b.unsetTop()
			}
		},

		func(s, b *Bitset, index uint) {
			s.set[index] |= b.set[index]
		},
	)
}

// Intersection another bitset to current bitset
func (s *Bitset) Intersection(b *Bitset) *Bitset {
	return s.bitsetOp(
		b,

		func(s, b *Bitset, length *uint) {
			bl, l := s.Length(0), *length

			if bl < l {
				s.unsetTop()
				s.Length(l)
			} else if bl > l {
				s.setTop()
			}
		},

		func(s, b *Bitset, index uint) {
			s.set[index] &= b.set[index]
		},
	)
}

// Diff calculate difference between current and another bitset
func (s *Bitset) Diff(b *Bitset) *Bitset {
	return s.bitsetOp(
		b,

		func(s, b *Bitset, length *uint) {
			if *length > s.Length(0) {
				*length = s.Length(0)
			} else {
				b.unsetTop()
			}
		},

		func(s, b *Bitset, index uint) {
			s.set[index] &= ^b.set[index]
		},
	)
}

// bitsetOp is common operation for union, intersection, diff
func (s *Bitset) bitsetOp(b *Bitset,
	lenFn func(s, b *Bitset, len *uint),
	opFn func(s, b *Bitset, index uint)) *Bitset {

	length := b.Length(0)
	if b == nil || b.Length(0) == 0 {
		return s
	}

	lenFn(s, b, &length)
	for i, n := Uint0, unitCount(length); i < n; i++ {
		opFn(s, b, i)
	}

	return s
}

// extend check if it's necessery to extend bitset's data store
func (s *Bitset) extend(index uint) {
	if index >= s.Length(0) {
		s.Length(index + 1)
	}
}

// unsetTop set bitset's top non-used units to 0
func (s *Bitset) unsetTop() {
	c := unitCount(s.length)
	for i := s.UnitCount() - 1; i >= c; i-- {
		s.set[i] = 0
	}
	s.set[c-1] &= (unitMax >> (c*unitLen - s.length))
}

// setTop set bitset's top non-used unitCount to 1
func (s *Bitset) setTop() {
	c := unitCount(s.length)
	for i := s.UnitCount() - 1; i >= c; i-- {
		s.set[i] = 1
	}
	s.set[c-1] |= (unitMax << (s.length - (c-1)*unitLen))
}

// Length change the bitset's length, if zero, only return current length.
//
// Only when new length is larger or half less than now, the allocation will occurred
func (s *Bitset) Length(l uint) uint {
	if l == 0 {
		return s.length
	}

	new := unitCount(l)
	if new > uint(cap(s.set)) || new <= s.UnitCount()>>1 {
		newSet := newUnits(new)
		copy(newSet, s.set)
		s.set = newSet
	} else {
		s.set = s.set[:new]
	}
	s.length = l

	return l
}

// UnitCount return bitset's unit count
func (s *Bitset) UnitCount() uint {
	return uint(len(s.set))
}

// UnitLen return unit length of bitset
func (s *Bitset) UnitLen() uint {
	return unitLen
}

// Uint return first uint unit
func (s *Bitset) Uint() uint {
	return uint(s.set[0])
}

// Uint64 return first uint64 unit
func (s *Bitset) Uint64() uint64 {
	return s.set[0]
}

// Clone return a new bitset same as current
func (s *Bitset) Clone() *Bitset {
	new := NewBitset(s.Length(0))
	copy(new.set, s.set)

	return new
}

// Bits return all index of bits set to 1
func (s *Bitset) Bits() []uint {
	res := make([]uint, s.BitCount())
	index := 0
	for i, l := Uint0, s.length; i < l; i++ {
		if s.IsSet(i) {
			res[index] = i
			index++
		}
	}

	return res
}

// BitCount return 1 bits count in bitset
func (s *Bitset) BitCount() int {
	var n int
	s.unsetTop()
	for i, l := 0, len(s.set); i < l; i++ {
		n += BitCount(s.set[i])
	}

	return n
}
