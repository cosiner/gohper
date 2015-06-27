package slices

func FitCapToLenForInt(slice []int) []int {
	if l := len(slice); l != cap(slice) {
		newslice := make([]int, l)
		copy(newslice, slice)
		return newslice
	}

	return slice
}

func IncrAppendForInt(slice []int, s int) []int {
	l := len(slice)
	newslice := make([]int, l+1)
	copy(newslice, slice)
	newslice[l] = s

	return newslice
}
