package slices

func FitCapToLenForInt32(slice []int32) []int32 {
	if l := len(slice); l != cap(slice) {
		newslice := make([]int32, l)
		copy(newslice, slice)
		return newslice
	}

	return slice
}

func IncrAppendForInt32(slice []int32, s int32) []int32 {
	l := len(slice)
	newslice := make([]int32, l+1)
	copy(newslice, slice)
	newslice[l] = s

	return newslice
}
