package types

// LightBitSet is a lightweight bitset implementation based on a uint64 number
// so each function which accept parameter index should gurentee that index is less
// than 64 for LightBitSet has not does that
type LightBitSet uint64

// NewLightBitSet create a new LightBitSet
func NewLightBitSet() *LightBitSet {
	var l LightBitSet
	return &l
}

// NewLightBitSetFrom create a new light bitset from given bits
func NewLightBitSetFrom(bits ...uint) *LightBitSet {
	var l uint64
	for _, b := range bits {
		l |= uint64(1 << b)
	}
	return ((*LightBitSet)(&l))
}

// ConvertToLightBitSet use a number as bitset
func ConvertToLightBitSet(b uint) *LightBitSet {
	var l LightBitSet = LightBitSet(b)
	return &l
}

// Set set bit at given index to 1
func (l *LightBitSet) Set(index uint) {
	*l |= 1 << index
}

// Unset set bit at given index to 0
func (l *LightBitSet) Unset(index uint) {
	*l &= ^(1 << index)
}

// IsSet chech whether bit at given index is set to 1
func (l *LightBitSet) IsSet(index uint) bool {
	return *l&(1<<index) != 0
}

// SetTo set bit at given index to 1 if val is true, else set to 0
func (l *LightBitSet) SetTo(index uint, val bool) {
	if val {
		l.Set(index)
	} else {
		l.Unset(index)
	}
}

// Flip flip bit at given index
func (l *LightBitSet) Flip(index uint) {
	*l ^= (1 << index)
}

// SetAll set all bits to 1
func (l *LightBitSet) SetAll() {
	*l = 1<<64 - 1
}

// SetAllBefore set all bits before index to 1, index bit is not included
func (l *LightBitSet) SetAllBefore(index uint) {
	*l |= (1<<index - 1)
}

func (l *LightBitSet) SetAllSince(index uint) {
	*l |= ^(1<<index - 1)
}

func (l *LightBitSet) UnsetAllBefore(index uint) {
	*l &= ^(1<<index - 1)
}

func (l *LightBitSet) UnsetAllSince(index uint) {
	*l &= (1<<index - 1)
}

// UnsetAll set all bits to 0
func (l *LightBitSet) UnsetAll() {
	*l = 0
}

// FlipAll flip all bits
func (l *LightBitSet) FlipAll() {
	*l = ^(*l)
}

// Uint return uint display of LightBitSet
func (l *LightBitSet) Uint() uint {
	return uint(*l)
}

// Uint64 return uint64 display of LightBitSet
func (l *LightBitSet) Uint64() uint64 {
	return uint64(*l)
}

// BitCount return the count of bits set to 1
func (l *LightBitSet) BitCount() int {
	return BitCount(uint64(*l))
}
