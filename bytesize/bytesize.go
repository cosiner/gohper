package bytesize

import (
	"bytes"
	"strconv"
)

const (
	// *BASE defines the base size of *B, include K,M,G,T,P
	BASE uint64 = 1 << 10 // 1024
	KB          = BASE
	MB          = BASE * KB
	GB          = BASE * MB
	TB          = BASE * GB
	PB          = BASE * TB
)

// Size convert byte count string to integer
// such as: 1K/k -> 1024, 1M/m -> 1024*1024
func Size(size string) (uint64, error) {
	var base uint64 = 1

	s := bytes.TrimSpace([]byte(size))
	switch s[len(s)-1] {
	case 'K', 'k':
		base *= KB
	case 'M', 'm':
		base *= MB
	case 'G', 'g':
		base *= GB
	case 'T', 't':
		base *= TB
	case 'P', 'p':
		base *= PB
	}

	if base > 1 {
		s = s[:len(s)-1]
	}

	bs, err := strconv.Atoi(string(s))
	if err != nil {
		return 0, err
	}

	if bs < 0 {
		bs = 0
	}

	return uint64(bs) * base, nil
}

// SizeDef is same as Size, on error return default size
func SizeDef(size string, defSize uint64) (s uint64) {
	var err error
	if s, err = Size(size); err != nil {
		return defSize
	}

	return s
}

// MustSize is same as Size, on error panic
func MustSize(size string) (s uint64) {
	s, err := Size(size)
	if err != nil {
		panic(err)
	}

	return s
}
