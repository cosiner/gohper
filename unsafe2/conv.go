package unsafe2

import (
	"reflect"
	"unsafe"
)

var Enable = true

// String bring a no copy convert from byte slice to string
// consider the risk
func String(b []byte) (s string) {
	if !Enable {
		return string(b)
	}

	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))

	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len

	return
}

// Bytes bring a no copy convert from string to byte slice
// consider the risk
func Bytes(s string) []byte {
	if !Enable {
		return []byte(s)
	}
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return BytesFromPtr(pstring.Data, pstring.Len)
}

func BytesFromPtr(ptr uintptr, len int) []byte {
	var b []byte
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pbytes.Data = uintptr(ptr)
	pbytes.Len = len
	pbytes.Cap = len
	return b
}
