package types

// Pow return power of a base number
func Pow(base, power uint) uint64 {
	n := uint64(1)
	ubase := uint64(base)
	for power > 0 {
		n *= ubase
		power--
	}

	return uint64(n)
}
