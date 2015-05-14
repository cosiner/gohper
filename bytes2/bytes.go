// Package bytes2 provide some helpful functions and multiple bytes pools.
package bytes2

import (
	"bytes"

	"github.com/cosiner/gohper/index"
)

// TrimSplit split bytes array, and trim space on each section
func TrimSplit(s, sep []byte) [][]byte {
	sp := bytes.Split(s, sep)
	for i, n := 0, len(sp); i < n; i++ {
		sp[i] = bytes.TrimSpace(sp[i])
	}

	return sp
}

// TrimAfter remove  bytes after delimeter, and trim space on remains
func TrimAfter(s []byte, delim []byte) []byte {
	if idx := bytes.Index(s, delim); idx >= 0 {
		s = s[:idx]
	}

	return bytes.TrimSpace(s)
}

// IsAllBytesIn check whether all bytes is in given encoding bytes
func IsAllBytesIn(bs []byte, encoding []byte) bool {
	var is = true
	for i := 0; i < len(bs) && is; i++ {
		is = index.ByteIn(bs[i], encoding...) >= 0
	}

	return is
}
