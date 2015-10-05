// Package bytes2 provide some helpful functions and multiple bytes pools.
package bytes2

import (
	"bytes"

	"github.com/cosiner/gohper/index"
)

// SplitAndTrim split bytes array, and trim space on each section
func SplitAndTrim(s, sep []byte) [][]byte {
	sp := bytes.Split(s, sep)
	for i, n := 0, len(sp); i < n; i++ {
		sp[i] = bytes.TrimSpace(sp[i])
	}

	return sp
}

// TrimAfter remove  bytes after delimeter, and trim space on remains
func TrimAfter(s, delim []byte) []byte {
	if idx := bytes.Index(s, delim); idx >= 0 {
		s = s[:idx]
	}

	return bytes.TrimSpace(s)
}

func TrimBefore(s, delim []byte) []byte {
	if idx := bytes.Index(s, delim); idx >= 0 {
		s = s[idx+len(delim):]
	}

	return bytes.TrimSpace(s)
}

// IsAllBytesIn check whether all bytes is in given encoding bytes
func IsAllBytesIn(bs, encoding []byte) bool {
	var is = true
	for i := 0; i < len(bs) && is; i++ {
		is = index.ByteIn(bs[i], encoding...) >= 0
	}

	return is
}

func MultipleLineOperate(s, delim []byte, operate func(line, delim []byte) []byte) []byte {
	var NEWLINE = []byte("\n")
	lines := bytes.Split(s, NEWLINE)
	for i := len(lines) - 1; i >= 0; i-- {
		lines[i] = operate(lines[i], delim)
	}

	return bytes.Join(lines, NEWLINE)
}

func TrimLastN(s, delim []byte, n int) []byte {
	if n <= 0 {
		n = -1
	}
	sl, dl := len(s), len(delim)
	for n != 0 && bytes.HasSuffix(s, delim) {
		s = s[:sl-dl]
		sl = len(s)
		n--
	}
	return s
}

func TrimFirstN(s, delim []byte, n int) []byte {
	if n <= 0 {
		n = -1
	}
	dl := len(delim)
	for n != 0 && bytes.HasPrefix(s, delim) {
		s = s[dl:]
		n--
	}
	return s
}

func LastIndexByte(bytes []byte, b byte) int {
	for i := len(bytes)-1; i >= 0; i-- {
		if bytes[i] == b {
			return i
		}
	}
	return -1
}