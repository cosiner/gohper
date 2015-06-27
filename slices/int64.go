package slices

func FitCapToLenForInt64(slice []int64) []int64 {
	if l := len(slice); l != cap(slice) {
		newslice := make([]int64, l)
		copy(newslice, slice)
		return newslice
	}

	return slice
}

func IncrAppendForInt64(slice []int64, s int64) []int64 {
	l := len(slice)
	newslice := make([]int64, l+1)
	copy(newslice, slice)
	newslice[l] = s

	return newslice
}
