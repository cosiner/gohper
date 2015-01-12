package types

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	. "github.com/cosiner/golib/errors"
)

// UnsafeString bring a no copy convert from byte slice to string
// consider the risk
func UnsafeString(b []byte) (s string) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len
	return
}

// UnsafeBytes bring a no copy convert from string to byte slice
// consider the risk
func UnsafeBytes(s string) (b []byte) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pbytes.Data = pstring.Data
	pbytes.Len = pstring.Len
	pbytes.Cap = pstring.Len
	return
}

// Str2Bool convert string to bool, error when not string
func Str2Bool(s string) (val bool, err error) {
	s = TrimLower(s)
	if s == "true" {
		val = true
	} else if s == "false" {
		val = false
	} else {
		err = Errorf("%s is not an bool string", s)
	}
	return
}

// MustStr2Bool is same as Str2Bool, on error it will panic
func MustStr2Bool(s string) bool {
	val, err := Str2Bool(s)
	if err != nil {
		panic(err)
	}
	return val
}

// Str2BoolDef convert string to bool, on error return default value
func Str2BoolDef(s string, def bool) bool {
	val, err := Str2Bool(s)
	if err != nil {
		val = def
	}
	return val
}

// Str2Int convert string to int, it's only a wrapper of strconv.Atoi
func Str2Int(s string) (int, error) {
	return strconv.Atoi(s)
}

// MustStr2Int convert string to int, on error panic
func MustStr2Int(s string) (val int) {
	val, err := Str2Int(s)
	OnErrPanic(err)
	return
}

// Str2IntDef convert string to int, on error use default value
func Str2IntDef(s string, def int) int {
	val, err := Str2Int(s)
	if err != nil {
		val = def
	}
	return val
}

// Int2Str convert integer to string use fmt.Sprintf
func Int2Str(val int) string {
	return fmt.Sprintf("%d", val)
}

// HexStr2Uint convert a hexadecimal string to uint
// if string is invalid, return an error
func HexStr2Uint(str string) (n uint, err error) {
	err = Errorf("Invalid hexadecimal string %s", str)
	if len(str) <= 2 {
		return
	} else {
		if head := str[:2]; head != "0x" && head != "0X" {
			return
		}
	}
	str = str[2:]
	for _, c := range str {
		if c >= '0' && c <= '9' {
			c = c - '0'
		} else if c >= 'a' && c <= 'f' {
			c = c - 'a' + 10
		} else if c >= 'A' && c <= 'F' {
			c = c - 'A' + 10
		} else {
			return
		}
		n = n << 4
		n |= uint(c)
	}
	err = nil
	return
}

// BytesToHexStr transfer binary to hex string
func BytesToHexStr(src []byte) string {
	return hex.EncodeToString(src)
}

// BytesToHex transfer binary to hex bytes
func BytesToHex(src []byte) []byte {
	dst := make([]byte, 0, 2*len(src))
	hex.Encode(dst, src)
	return dst
}

// HexToBytes transfer hex bytes to binary
func HexToBytes(src []byte) []byte {
	dst := make([]byte, 0, len(src)/2)
	hex.Decode(dst, src)
	return dst
}

// ReverseBits reverse all bits in number
func ReverseBits(num uint) uint {
	var n uint
	size := uint(unsafe.Sizeof(num))
	for s := size * 8; s > 0; s-- {
		n = n << 1
		n |= (num & 1)
		num = num >> 1
	}
	return n
}

// ReverseByte reverse all bits for a byte
func ReverseByte(num uint8) uint8 {
	var n uint8
	for s := 8; s > 0; s-- {
		n = n << 1
		n |= (num & 1)
		num = num >> 1
	}
	return n
}
