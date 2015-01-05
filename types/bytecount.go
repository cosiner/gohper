package types

import (
	"bytes"
	"github.com/cosiner/golib/errors"
	"strconv"
)

const (
	// *BYTE_BASE defines the base size of *B, include K,M,G,T,P
	BYTE_BASE  = 1024
	KBYTE_BASE = BYTE_BASE * 1
	MBYTE_BASE = BYTE_BASE * KBYTE_BASE
	GBYTE_BASE = BYTE_BASE * MBYTE_BASE
	TBYTE_BASE = BYTE_BASE * GBYTE_BASE
	PBYTE_BASE = BYTE_BASE * TBYTE_BASE
)

// Str2Bytes convert byte count string to integer
// such as: 1K/k -> 1024, 1M/m -> 1024*1024
func Str2Bytes(size string) (uint64, error) {
	var base uint64 = 1
	s := []byte(size)
	s = bytes.TrimSpace(s)
	switch size[len(size)-1] {
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

// Str2BytesDef is same as Str2Bytes, on error return default size
func Str2BytesDef(size string, defSize uint64) (s uint64) {
	var err error
	if s, err = Str2Bytes(size); err != nil {
		return defSize
	}
	return s
}

// MustStr2Bytes is same as Str2Bytes, on error panic
func MustStr2Bytes(size string) (s uint64) {
	errors.ErrorPanic(func() error {
		var err error
		s, err = Str2Bytes(size)
		return err
	})
	return
}
