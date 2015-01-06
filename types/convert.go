package types

import (
	. "github.com/cosiner/golib/errors"
	"reflect"
	"strconv"
	"unsafe"
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
	ErrorPanic(func() (err error) {
		val, err = Str2Int(s)
		return
	})
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
