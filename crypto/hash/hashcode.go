package hashcode

func HashCode(str []byte, mod uint) uint {
	hash := BKDR(str)
	if mod != 0 {
		hash &= mod - 1
	}

	return hash
}

func SDBM(str []byte) uint {
	var hash uint

	for _, b := range str {
		hash = uint(b) + (hash << 6) + (hash << 16) - hash
	}

	return hash & 0x7FFFFFFF
}

func RS(str []byte) uint {
	var c uint = 378551
	var a uint = 63689
	var hash uint

	for _, b := range str {
		hash = hash*a + uint(b)
		a *= c
	}

	return hash & 0x7FFFFFFF
}

func JS(str []byte) uint {
	var hash uint = 1315423911

	for _, b := range str {
		hash ^= ((hash << 5) + uint(b) + (hash >> 2))
	}

	return hash & 0x7FFFFFFF
}

func ELF(str []byte) uint {
	var hash uint
	var x uint

	for _, b := range str {
		hash = (hash << 4) + uint(b)
		if x = hash & 0xF0000000; x != 0 {
			hash ^= (x >> 24)
			hash &= ^x
		}
	}

	return hash & 0x7FFFFFFF
}

func BKDR(str []byte) uint {
	var seed uint = 131 // 31 131 1313 13131 131313 etc..
	var hash uint

	for _, b := range str {
		hash = hash*seed + uint(b)
	}

	return hash & 0x7FFFFFFF
}

func DJB(str []byte) uint {
	var hash uint = 5381

	for _, b := range str {
		hash += (hash << 5) + uint(b)
	}

	return hash & 0x7FFFFFFF
}

func AP(str []byte) uint {
	var hash uint
	var i uint

	for _, b := range str {
		if i&1 == 0 {
			hash ^= ((hash << 7) ^ uint(b) ^ (hash >> 3))
		} else {
			hash ^= (^((hash << 11) ^ uint(b) ^ (hash >> 5)))
		}
	}

	return hash & 0x7FFFFFFF
}
