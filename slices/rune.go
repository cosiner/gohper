package slices

func FitCapToLenForRune(slice []rune) []rune {
	if l := len(slice); l != cap(slice) {
		newslice := make([]rune, l)
		copy(newslice, slice)
		return newslice
	}

	return slice
}

func IncrAppendForRune(slice []rune, s rune) []rune {
	l := len(slice)
	newslice := make([]rune, l+1)
	copy(newslice, slice)
	newslice[l] = s

	return newslice
}
