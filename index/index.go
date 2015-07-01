package index

// BitIn test whether the bit at index is set to 1, if true, return 1 << index, else 0
func BitIn(index int, bitset uint) uint {
	if index < 0 {
		return 0
	}

	var idx uint = 1 << uint(index)
	if idx&bitset != 0 {
		return idx
	}

	return 0
}

// BitNotIn test whether the bit at index is set to 0, if true, return 1 << index, else 0
func BitNotIn(index int, bitset uint) uint {
	if index < 0 {
		return 0
	}

	var idx uint = 1 << uint(index)
	if idx&bitset == 0 {
		return idx
	}

	return 0
}

// RuneIn return the index that rune in rune list or -1 if not exist
func RuneIn(ru rune, rs ...rune) int {
	for index, r := range rs {
		if r == ru {
			return index
		}
	}

	return -1
}

// ByteIn return the index that byte in byte list or -1 if not exist
func ByteIn(b byte, bs ...byte) int {
	for index, c := range bs {
		if b == c {
			return index
		}
	}

	return -1
}

// StringIn return the index of string to find in a string slice or -1 if not found
func StringIn(str string, strs []string) int {
	for i, s := range strs {
		if s == str {
			return i
		}
	}

	return -1
}

// SortedNumberIn search index of number in sorted numbers with ascending order
func SortedNumberIn(n int, nums ...int) int {
	for l, h := 0, len(nums)-1; l <= h; {
		m := l + (h-l)>>1

		if c := nums[m]; c == n {
			return m
		} else if c < n {
			l = m + 1
		} else {
			h = m - 1
		}
	}

	return -1
}

// CharIn search index of char in sorted string with ascending order
func CharIn(b byte, s string) int {
	for l, h := 0, len(s)-1; l <= h; {
		m := l + (h-l)>>1

		if c := s[m]; c == b {
			return m
		} else if c < b {
			l = m + 1
		} else {
			h = m - 1
		}
	}

	return -1
}
