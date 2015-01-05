// Package types implements some type relevant tools
// bitset, string, bytes
// M/N  = M >> n : N= 2 ** n
// M%N = M & (N-1):N = 2 ** n
package types

// u_1 is uint 1
const (
	uint0          uint   = 0
	unitLenLogN           = 6 // log64 = 6
	unitLen               = 1 << unitLenLogN
	unitMax        uint64 = 1<<unitLen - 1
	minShrinkCount        = 16
)

// BitSet is a bitset
type BitSet struct {
	count       uint // bitset unit count
	shrinkCount uint // bitset shrink unit count
	set         []uint64
}

// unitCount return unit count need for the count
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

// NewBitSet return a new bitset with gived count
// Notice: the parameter count represent 64bits
func NewBitSet(count uint) *BitSet {
	return &BitSet{count, minShrinkCount, make([]uint64, count)}
}

// Length return bitset count
func (bs *BitSet) Length() uint {
	return bs.count * unitLen
}

// Clone return a new bitset same as current
func (bs *BitSet) Clone() *BitSet {
	newBitSet := NewBitSet(bs.count)
	if bs.count > 0 {
		copy(newBitSet.set, bs.set)
	}
	return newBitSet
}

// Shrink shrink the bitset, before call this, must setup minshrinkcount,
// otherwise, may make no difference
func (bs *BitSet) Shrink(count uint) {
	bs.changeUnitCount(count)
}

// changeUnitCount change the bitset's count
func (bs *BitSet) changeUnitCount(count uint) {
	oldCount := uint(len(bs.set))
	if oldCount < count ||
		(count*2 <= oldCount &&
			oldCount > bs.shrinkCount) {
		// if unit count is not enough or new unit count is twice small than old
		// and old unit count is too large, then make a new set and copy data
		newSet := make([]uint64, count)
		if bs.set != nil {
			copy(bs.set, newSet)
			bs.set = nil
		}
		bs.set = newSet
	}
	bs.count = count
}

// ChangeShrinkCount change set shrink count, only used when bitset is shrinking
func (bs *BitSet) ChangeShrinkCount(count uint) {
	if count > 1 {
		bs.shrinkCount = count
	}
}

// Set set index bit to 1
// if index large then bitset count, will expand the bitset
func (bs *BitSet) Set(index uint) *BitSet {
	if index >= bs.count*unitLen {
		bs.changeUnitCount(unitCount(index + 1))
	}
	bs.set[unitPos(index)] |= 1 << unitIndex(index)
	return bs
}

// Set set index bit to 1
func (bs *BitSet) SetAll() *BitSet {
	return bs.unitOp(func(index uint) {
		bs.set[index] = unitMax
	})
}

// Unset set index bit to 0
func (bs *BitSet) UnSet(index uint) *BitSet {
	if index < bs.count*unitLen {
		bs.set[unitPos(index)] &= ^(1 << unitIndex(index))
	}
	return bs
}

// UnSetAll set all bits to 0
func (bs *BitSet) UnSetAll() *BitSet {
	return bs.unitOp(func(index uint) {
		bs.set[index] = 0
	})
}

// Flip flip the index bit
func (bs *BitSet) Flip(index uint) *BitSet {
	if index >= bs.count {
		return bs.Set(index)
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

// IsSet check whether or not index bit is set
func (bs *BitSet) IsSet(index uint) bool {
	return index < bs.count && (bs.set[unitPos(index)]&(1<<unitIndex(index))) != 0
}

// SetTo set index bit to 1 if value is true, otherwise 0
func (bs *BitSet) SetTo(index uint, value bool) *BitSet {
	if value {
		return bs.Set(index)
	}
	return bs.UnSet(index)
}

func (bs *BitSet) BitCount() uint {
	var n uint = 0
	bs.unitOp(func(index uint) {
		n += bitCount(bs.set[index])
	})
	return n
}

// Union union another bitset to current bitset
// if want union to a new bitset instead of change current bitset,
// please call Clone() first to create a new bitset, then call Union
// on new bitset
func (bs *BitSet) Union(b *BitSet) *BitSet {
	count := b.count
	if b == nil || count == 0 {
		return bs
	}
	if bs.count < count {
		bs.changeUnitCount(count)
	}
	for i := uint0; i < count; i++ {
		bs.set[i] |= b.set[i]
	}

	return bs
}

// Intersection intersection another bitset to current bitset
func (bs *BitSet) Intersection(b *BitSet) *BitSet {
	count := b.count
	if b == nil || count == 0 {
		return bs
	}
	bsCount := bs.count
	if bsCount < count {
		bs.changeUnitCount(count)
	}
	for i := uint0; i < count; i++ {
		bs.set[i] &= b.set[i]
	}
	for i := count; i < bsCount; i++ {
		bs.set[i] = 0
	}
	return bs
}

// Diff calculate difference between current and another bitset
func (bs *BitSet) Diff(b *BitSet) *BitSet {
	count := b.count
	if b == nil || count == 0 {
		return bs
	}
	if count > bs.count {
		count = bs.count
	}
	for i := uint0; i < count; i++ {
		bs.set[i] &= ^b.set[i]
	}
	return bs
}

// unitOp iter the bitset unit, apply function to each unit
func (bs *BitSet) unitOp(f func(index uint)) *BitSet {
	for i, n := uint0, bs.count; i < n; i++ {
		f(i)
	}
	return bs
}

// count of 1 bit
func bitCount(n uint64) uint {
	n -= (n >> 1) & 0x5555555555555555
	n = (n>>2)&0x3333333333333333 + n&0x3333333333333333
	n += n >> 4
	n &= 0x0f0f0f0f0f0f0f0f
	n *= 0x0101010101010101
	return uint(n >> 56)
}

// In 测试第index位在bitset中是否被设置,如果设置了,返回index位被设置的bitset 即 1 << index
func In(index int, bitset uint) (i uint) {
	if index >= 0 {
		var idx uint = 1 << uint(index)
		if idx&bitset != 0 {
			i = idx
		}
	}
	return
}

// NotIn 测试第index位在bitset中是否被设置,如果没设置, 返回index位被设置的bitset 即 1 << index
func NotIn(index int, bitset uint) (i uint) {
	if index >= 0 {
		var idx uint = 1 << uint(index)
		if idx&bitset == 0 {
			i = idx
		}
	}
	return
}
