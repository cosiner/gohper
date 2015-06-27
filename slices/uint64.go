package slices

func FitCapToLenForUint64(slice []uint64) []uint64 {
	if l := len(slice); l != cap(slice) {
		newslice := make([]uint64, l)
		copy(newslice, slice)
		return newslice
	}

	return slice
}

func IncrAppendForUint64(slice []uint64, s uint64) []uint64 {
	l := len(slice)
	newslice := make([]uint64, l+1)
	copy(newslice, slice)
	newslice[l] = s

	return newslice
}
