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
