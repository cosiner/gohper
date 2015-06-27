package slices

func FitCapToLenForByte(slice []byte) []byte {
	if l := len(slice); l != cap(slice) {
		newslice := make([]byte, l)
		copy(newslice, slice)
		return newslice
	}

	return slice
}

func IncrAppendForByte(slice []byte, s byte) []byte {
	l := len(slice)
	newslice := make([]byte, l+1)
	copy(newslice, slice)
	newslice[l] = s

	return newslice
}
