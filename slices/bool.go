package slices

func FitCapToLenForBool(slice []bool) []bool {
	if l := len(slice); l != cap(slice) {
		newslice := make([]bool, l)
		copy(newslice, slice)
		return newslice
	}

	return slice
}

func IncrAppendForBool(slice []bool, s bool) []bool {
	l := len(slice)
	newslice := make([]bool, l+1)
	copy(newslice, slice)
	newslice[l] = s

	return newslice
}
