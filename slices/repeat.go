package slices

func RepeatBytes(b byte, n int) []byte {
	bs := make([]byte, n)
	for n = n - 1; n >= 0; n-- {
		bs[n] = b
	}

	return bs
}

func RepeatStrings(s string, n int) []string {
	bs := make([]string, n)
	for n = n - 1; n >= 0; n-- {
		bs[n] = s
	}

	return bs
}

func RepeatInts(i int, n int) []int {
	bs := make([]int, n)
	for n = n - 1; n >= 0; n-- {
		bs[n] = i
	}

	return bs
}
