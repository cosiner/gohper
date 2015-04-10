package types

import (
	"bytes"
	"strconv"

	. "github.com/cosiner/gohper/lib/errors"
)

const (
	// *BYTE_BASE defines the base size of *B, include K,M,G,T,P
	BYTE_BASE  = 1 << 10
	KBYTE_BASE = BYTE_BASE
	MBYTE_BASE = BYTE_BASE * KBYTE_BASE
	GBYTE_BASE = BYTE_BASE * MBYTE_BASE
	TBYTE_BASE = BYTE_BASE * GBYTE_BASE
	PBYTE_BASE = BYTE_BASE * TBYTE_BASE
)

// BytesCount convert byte count string to integer
// such as: 1K/k -> 1024, 1M/m -> 1024*1024
func BytesCount(size string) (uint64, error) {
	var base uint64 = 1
	s := bytes.TrimSpace([]byte(size))
	switch s[len(s)-1] {
	case 'K', 'k':
		base *= KBYTE_BASE
	case 'M', 'm':
		base *= MBYTE_BASE
	case 'G', 'g':
		base *= GBYTE_BASE
	case 'T', 't':
		base *= TBYTE_BASE
	case 'P', 'p':
		base *= PBYTE_BASE
	}
	if base > 1 {
		s = s[:len(s)-1]
	}
	bs, err := strconv.Atoi(string(s))
	if err != nil {
		return 0, err
	}
	return uint64(bs) * base, nil
}

// BytesCountDef is same as BytesCount, on error return default size
func BytesCountDef(size string, defSize uint64) (s uint64) {
	var err error
	if s, err = BytesCount(size); err != nil {
		return defSize
	}
	return s
}

// MustBytesCount is same as BytesCount, on error panic
func MustBytesCount(size string) (s uint64) {
	s, err := BytesCount(size)
	OnErrPanic(err)
	return s
}

// TODO
func Bytes2Str(size uint64) string {
	return ""
}
