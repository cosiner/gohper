package slices

func FitCapToLenForString(slice []string) []string {
	if l := len(slice); l != cap(slice) {
		newslice := make([]string, l)
		copy(newslice, slice)
		return newslice
	}

	return slice
}

func IncrAppendForString(slice []string, s string) []string {
	l := len(slice)
	newslice := make([]string, l+1)
	copy(newslice, slice)
	newslice[l] = s

	return newslice
}
