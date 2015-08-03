package slices

func RemoveElement(slice []interface{}, index int) []interface{} {
	size := len(slice)
	if index >= size {
		return slice
	}

	if size >= 2*(index+1) {
		for ; index > 0; index-- {
			slice[index] = slice[index-1]
		}
		slice = slice[1:]
	} else {
		for ; index < size-1; index++ {
			slice[index] = slice[index+1]
		}
		slice = slice[:size-1]
	}

	return slice
}

func InterfacesMap(eles []interface{}, mapper func(interface{}) interface{}) []interface{} {
	for i, l := 0, len(eles); i < l; i++ {
		eles[i] = mapper(eles[i])
	}
	return eles
}

func InterfacesFilter(eles []interface{}, filter func(interface{}) bool) []interface{} {
	var newEles []interface{}
	for i, l := 0, len(eles); i < l; i++ {
		if e := eles[i]; filter(e) {
			newEles = append(newEles, e)
		}
	}

	return newEles
}

func InterfacesFilterInplace(eles []interface{}, filter func(interface{}) bool) []interface{} {
	var prev = -1

	for i, l := 0, len(eles); i < l; i++ {
		if e := eles[i]; filter(e) {
			prev++
			eles[prev] = e
		}
	}

	return eles[:prev+1]
}
