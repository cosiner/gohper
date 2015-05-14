package bitset

// Bits is a lightweight bitset implementation based on a uint64 number,
// it's not safety for concurrent.
type Bits uint64

func NewBits() *Bits {
	var s Bits

	return &s
}

// BitsList create a Bits, set all bits in list to 1
func BitsList(bits ...uint) *Bits {
	var s uint64
	for _, b := range bits {
		s |= uint64(1 << b)
	}

	return ((*Bits)(&s))
}

// BitsFrom create a Bits, all bits of number is copied
func BitsFrom(b uint) *Bits {
	var s Bits = Bits(b)

	return &s
}

// Set bit at given index to 1
func (s *Bits) Set(index uint) {
	*s |= 1 << index
}

// SetAll set all bits to 1
func (s *Bits) SetAll() {
	*s = 1<<64 - 1
}

// Unset set bit at given index to 0
func (s *Bits) Unset(index uint) {
	*s &= ^(1 << index)
}

// UnsetAll set all bits to 0
func (s *Bits) UnsetAll() {
	*s = 0
}

// IsSet chech whether bit at given index is set to 1
func (s *Bits) IsSet(index uint) bool {
	return *s&(1<<index) != 0
}

// SetTo set bit at given index to 1 if val is true, else set to 0
func (s *Bits) SetTo(index uint, val bool) {
	if val {
		s.Set(index)
	} else {
		s.Unset(index)
	}
}

// Flip bit at given index
func (s *Bits) Flip(index uint) {
	*s ^= (1 << index)
}

// FlipAll flip all bits
func (s *Bits) FlipAll() {
	*s = ^(*s)
}

// SetBefore set all bits before index to 1, index bit is not included
func (s *Bits) SetBefore(index uint) {
	*s |= (1<<index - 1)
}

// SetSince set all bits since index to 1, index bit is include
func (s *Bits) SetSince(index uint) {
	*s |= ^(1<<index - 1)
}

// UnsetBefore set all bits before index to 0, index bit is not included
func (s *Bits) UnsetBefore(index uint) {
	*s &= ^(1<<index - 1)
}

// UnsetSince set all bits since index to 0, index bit is include
func (s *Bits) UnsetSince(index uint) {
	*s &= (1<<index - 1)
}

// Uint return uint display of Bits
func (s *Bits) Uint() uint {
	return uint(*s)
}

// Uint64 return uint64 display of Bits
func (s *Bits) Uint64() uint64 {
	return uint64(*s)
}

// BitCount return the count of bits set to 1
func (s *Bits) BitCount() int {
	return BitCount(uint64(*s))
}

// BitCount return count of 1 bit in uint64
func BitCount(n uint64) int {
	n -= (n >> 1) & 0x5555555555555555
	n = (n>>2)&0x3333333333333333 + n&0x3333333333333333
	n += n >> 4
	n &= 0x0f0f0f0f0f0f0f0f
	n *= 0x0101010101010101

	return int(n >> 56)
}

// BitCountUint return count of 1 bit in uint
func BitCountUint(x uint) int {
	var n = uint64(x)
	n -= (n >> 1) & 0x5555555555555555
	n = (n>>2)&0x3333333333333333 + n&0x3333333333333333
	n += n >> 4
	n &= 0x0f0f0f0f0f0f0f0f
	n *= 0x0101010101010101

	return int(n >> 56)
}
