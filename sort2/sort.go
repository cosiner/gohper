package sort2

import "sort"

type bytes []byte

func (b bytes) Len() int {
	return len(b)
}

func (b bytes) Less(i, j int) bool {
	return b[i] < b[j]
}

func (b bytes) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func Bytes(b []byte) []byte {
	sort.Sort(bytes(b))
	return b
}

func String(s string) string {
	return string(Bytes([]byte(s)))
}
