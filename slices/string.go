package slices

func FitCapToLenString(slice []string) []string {
	if l := len(slice); l != cap(slice) {
		newslice := make([]string, l)
		copy(newslice, slice)
		return newslice
	}

	return slice
}

func IncrAppendString(slice []string, s string) []string {
	l := len(slice)
	newslice := make([]string, l+1)
	copy(newslice, slice)
	newslice[l] = s

	return newslice
}

func StringsMap(eles []string, mapper func(string) string) []string {
	for i, l := 0, len(eles); i < l; i++ {
		eles[i] = mapper(eles[i])
	}
	return eles
}

func StringsFilter(eles []string, filter func(string) bool) []string {
	var newEles []string
	for i, l := 0, len(eles); i < l; i++ {
		if e := eles[i]; filter(e) {
			newEles = append(newEles, e)
		}
	}

	return newEles
}

func StringsFilterInplace(eles []string, filter func(string) bool) []string {
	var prev = -1

	for i, l := 0, len(eles); i < l; i++ {
		if e := eles[i]; filter(e) {
			prev++
			eles[prev] = e
		}
	}

	return eles[:prev+1]
}
