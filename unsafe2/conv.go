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
func Bytes(s string) (b []byte) {
	if !Enable {
		return []byte(s)
	}
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))

	pbytes.Data = pstring.Data
	pbytes.Len = pstring.Len
	pbytes.Cap = pstring.Len

	return
}
