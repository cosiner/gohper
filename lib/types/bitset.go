// Package types implements some type relevant tools
// bitset, string, bytes
package types

// u_1 is uint 1
const (
	Uint0       uint   = 0
	unitLenLogN        = 6 // log64 = 6
	unitLen            = 1 << unitLenLogN
	unitMax     uint64 = 1<<unitLen - 1
)

// BitSet is a bitset
type BitSet struct {
	length uint     // bitset length
	set    []uint64 // bitset data store
}

// NewBitSet return a new bitset with given length, all index params are set to 1
func NewBitSet(length uint, indexs ...uint) (bs *BitSet) {
	if length != 0 {
		bs = &BitSet{length, newUnitSet(unitCount(length))}
		for _, idx := range indexs {
			bs.Set(idx)
		}
	}
	return
}

// Len return bitset length
func (bs *BitSet) Len() uint {
	return bs.length
}

// Cap return bitset's capacity
func (bs *BitSet) Cap() uint {
	return bs.UnitLen() * bs.UnitCount()
}

// UnitCount return bitset's unit count
func (bs *BitSet) UnitCount() uint {
	return uint(len(bs.set))
}

// UnitLen return unit length of bitset
func (bs *BitSet) UnitLen() uint {
	return unitLen
}

// Clone return a new bitset same as current
func (bs *BitSet) Clone() *BitSet {
	newBitSet := NewBitSet(bs.Len())
	if bs.Len() > 0 {
		copy(newBitSet.set, bs.set)
	}
	return newBitSet
}

// Shrink apply bitset for shrink operation
func (bs *BitSet) Shrink() *BitSet {
	return bs.ChangeLength(bs.length)
}

// ChangeLength change the bitset's length
func (bs *BitSet) ChangeLength(length uint) *BitSet {
	newCount := unitCount(length)
	oldCount := bs.UnitCount()
	if oldCount < newCount || newCount*2 <= oldCount {
		newSet := newUnitSet(newCount)
		if bs.set != nil {
			copy(newSet, bs.set)
		}
		bs.set = newSet
	}
	bs.length = length
	return bs
}

// Set set index bit to 1
// if index large then bitset count, will expand the bitset
func (bs *BitSet) Set(index uint) *BitSet {
	bs.extend(index)
	bs.set[unitPos(index)] |= 1 << unitIndex(index)
	return bs
}

// SetAll set all bits to 1
func (bs *BitSet) SetAll() *BitSet {
	return bs.unitOp(func(index uint) {
		bs.set[index] = unitMax
	})
}

// Unset set index bit to 0
// if index is larger than bitset length, bitset must be extended
func (bs *BitSet) Unset(index uint) *BitSet {
	bs.extend(index)
	bs.set[unitPos(index)] &= ^(1 << unitIndex(index))
	return bs
}

// UnsetAll set all bits to 0
func (bs *BitSet) UnsetAll() *BitSet {
	return bs.unitOp(func(index uint) {
		bs.set[index] = 0
	})
}

// Flip flip the index bit
func (bs *BitSet) Flip(index uint) *BitSet {
	if index >= bs.Len() {
		bs.Set(index)
	}
	bs.set[unitPos(index)] ^= 1 << unitIndex(index)
	return bs
}

// FlipAll flip all the index bit
func (bs *BitSet) FlipAll() *BitSet {
	return bs.unitOp(func(index uint) {
		bs.set[index] = ^bs.set[index]
	})
}

// Except set all bits except given index to 1, the except bits set to 0
func (bs *BitSet) Except(index ...uint) *BitSet {
	bs.SetAll()
	for _, i := range index {
		bs.Unset(i)
	}
	return bs
}

// IsSet check whether or not index bit is set
func (bs *BitSet) IsSet(index uint) bool {
	return index < bs.Len() && (bs.set[unitPos(index)]&(1<<unitIndex(index))) != 0
}

// SetTo set index bit to 1 if value is true, otherwise 0
func (bs *BitSet) SetTo(index uint, value bool) *BitSet {
	if value {
		return bs.Set(index)
	}
	return bs.Unset(index)
}

// Bits return all index of bits set to 1
func (bs *BitSet) Bits() (res []uint) {
	res = make([]uint, 0, bs.BitCount())
	index := 0
	for i, l := Uint0, bs.Len(); i < l && bs.IsSet(i); i++ {
		res[index] = i
		index++
	}
	return
}

// BitCount count 1 bits
func (bs *BitSet) BitCount() int {
	var n int
	bs.clearTop()
	bs.unitOp(func(index uint) {
		n += BitCount(bs.set[index])
	})
	return n
}

// Uint return first uint unit
func (bs *BitSet) Uint() uint {
	return uint(bs.set[0])
}

// Uint64 return first uint64 unit
func (bs *BitSet) Uint64() uint64 {
	return bs.set[0]
}

// Union union another bitset to current bitset
// if want union to a new bitset instead of change current bitset,
// please call Clone() first to create a new bitset, then call Union
// on new bitset
func (bs *BitSet) Union(b *BitSet) *BitSet {
	return bs.bitsetOp(b,
		func(length *uint) {
			bl, l := bs.Len(), *length
			if bl < l {
				bs.clearTop()
				bs.ChangeLength(l)
			} else if bl > l {
				b.clearTop()
			}
		},
		func(index uint) {
			bs.set[index] |= b.set[index]
		})
}

// Intersection intersection another bitset to current bitset
func (bs *BitSet) Intersection(b *BitSet) *BitSet {
	return bs.bitsetOp(b,
		func(length *uint) {
			bl, l := bs.Len(), *length
			if bl < l {
				bs.clearTop()
				bs.ChangeLength(l)
			} else if bl > l {
				bs.setTop()
			}
		},
		func(index uint) {
			bs.set[index] &= b.set[index]
		})
}

// Diff calculate difference between current and another bitset
func (bs *BitSet) Diff(b *BitSet) *BitSet {
	return bs.bitsetOp(b,
		func(length *uint) {
			if *length > bs.Len() {
				*length = bs.Len()
			} else {
				b.clearTop()
			}
		},
		func(index uint) {
			bs.set[index] &= ^b.set[index]
		})
}

// bitsetOp is common operation for union, intersection, diff
func (bs *BitSet) bitsetOp(b *BitSet, lenFn func(*uint), opFn func(index uint)) *BitSet {
	length := b.Len()
	if b == nil || b.Len() == 0 {
		return bs
	}
	lenFn(&length)
	for i, n := Uint0, unitCount(length); i < n; i++ {
		opFn(i)
	}
	return bs
}

// extend check if it's necessery to extend bitset's data store
func (bs *BitSet) extend(index uint) {
	if index >= bs.Len() {
		bs.ChangeLength(index + 1)
	}
}

// clearTop clear bitset's top non-used unit, all these bits are set to zero
func (bs *BitSet) clearTop() {
	units := unitCount(bs.length)
	for i := bs.UnitCount() - 1; i >= units; i-- {
		bs.set[i] = 0
	}
	bs.set[units-1] &= (unitMax >> (units*unitLen - bs.length))
}

// setTop set bitset's top non-used unit,  to 1
func (bs *BitSet) setTop() {
	units := unitCount(bs.length)
	for i := bs.UnitCount() - 1; i >= units; i-- {
		bs.set[i] = 1
	}
	bs.set[units-1] |= (unitMax << (bs.length - (units-1)*unitLen))
}

// newUnitSet create a new unit set has given unit count for bitset
func newUnitSet(count uint) []uint64 {
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

// unitOp iter the bitset unit, apply function to each unit
func (bs *BitSet) unitOp(f func(index uint)) *BitSet {
	for i, n := Uint0, unitCount(bs.Len()); i < n; i++ {
		f(i)
	}
	return bs
}

// count of 1 bit
func BitCount(n uint64) int {
	n -= (n >> 1) & 0x5555555555555555
	n = (n>>2)&0x3333333333333333 + n&0x3333333333333333
	n += n >> 4
	n &= 0x0f0f0f0f0f0f0f0f
	n *= 0x0101010101010101
	return int(n >> 56)
}

func BitCountUint(x uint) int {
	var n = uint64(x)
	n -= (n >> 1) & 0x5555555555555555
	n = (n>>2)&0x3333333333333333 + n&0x3333333333333333
	n += n >> 4
	n &= 0x0f0f0f0f0f0f0f0f
	n *= 0x0101010101010101
	return int(n >> 56)
}
