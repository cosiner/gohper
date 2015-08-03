package slices

func FitCapToLenInt(slice []int) []int {
	if l := len(slice); l != cap(slice) {
		newslice := make([]int, l)
		copy(newslice, slice)
		return newslice
	}

	return slice
}

func IncrAppendInt(slice []int, s int) []int {
	l := len(slice)
	newslice := make([]int, l+1)
	copy(newslice, slice)
	newslice[l] = s

	return newslice
}

func IntsMap(eles []int, mapper func(int) int) []int {
	for i, l := 0, len(eles); i < l; i++ {
		eles[i] = mapper(eles[i])
	}
	return eles
}

func IntsFilter(eles []int, filter func(int) bool) []int {
	var newEles []int
	for i, l := 0, len(eles); i < l; i++ {
		if e := eles[i]; filter(e) {
			newEles = append(newEles, e)
		}
	}

	return newEles
}

func IntsFilterInplace(eles []int, filter func(int) bool) []int {
	var prev = -1

	for i, l := 0, len(eles); i < l; i++ {
		if e := eles[i]; filter(e) {
			prev++
			eles[prev] = e
		}
	}

	return eles[:prev+1]
}
