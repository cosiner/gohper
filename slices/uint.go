package slices

func FitCapToLenForUint(slice []uint) []uint {
	if l := len(slice); l != cap(slice) {
		newslice := make([]uint, l)
		copy(newslice, slice)
		return newslice
	}

	return slice
}

func IncrAppendForUint(slice []uint, s uint) []uint {
	l := len(slice)
	newslice := make([]uint, l+1)
	copy(newslice, slice)
	newslice[l] = s

	return newslice
}
