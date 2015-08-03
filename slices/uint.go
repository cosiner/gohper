package slices

func FitCapToLenUint(slice []uint) []uint {
	if l := len(slice); l != cap(slice) {
		newslice := make([]uint, l)
		copy(newslice, slice)
		return newslice
	}

	return slice
}

func IncrAppendUint(slice []uint, s uint) []uint {
	l := len(slice)
	newslice := make([]uint, l+1)
	copy(newslice, slice)
	newslice[l] = s

	return newslice
}

func UintsMap(eles []uint, mapper func(uint) uint) []uint {
	for i, l := 0, len(eles); i < l; i++ {
		eles[i] = mapper(eles[i])
	}
	return eles
}

func UintsFilter(eles []uint, filter func(uint) bool) []uint {
	var newEles []uint
	for i, l := 0, len(eles); i < l; i++ {
		if e := eles[i]; filter(e) {
			newEles = append(newEles, e)
		}
	}

	return newEles
}

func UintsFilterInplace(eles []uint, filter func(uint) bool) []uint {
	var prev = -1

	for i, l := 0, len(eles); i < l; i++ {
		if e := eles[i]; filter(e) {
			prev++
			eles[prev] = e
		}
	}

	return eles[:prev+1]
}
