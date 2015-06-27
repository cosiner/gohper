package slices

func FitCapToLenForUint32(slice []uint32) []uint32 {
	if l := len(slice); l != cap(slice) {
		newslice := make([]uint32, l)
		copy(newslice, slice)
		return newslice
	}

	return slice
}

func IncrAppendForUint32(slice []uint32, s uint32) []uint32 {
	l := len(slice)
	newslice := make([]uint32, l+1)
	copy(newslice, slice)
	newslice[l] = s

	return newslice
}
