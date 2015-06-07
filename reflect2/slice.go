package reflect2

import (
	"reflect"

	"github.com/cosiner/gohper/errors"
)

var ErrNonSlice = errors.Err("parameter is not a slice")

// IncrAppend will always create a new slice with the increased cap and length,
// then append the element, copy original elements and append new element to new slice,
// nil slice is allowed, but it must be a slice
func IncrAppend(s interface{}, e interface{}) interface{} {
	src := reflect.ValueOf(s)
	errors.Assert(src.Kind() == reflect.Slice, ErrNonSlice)

	len := src.Len()
	dst := reflect.MakeSlice(src.Type(), len+1, src.Cap()+1)
	reflect.Copy(dst, src)
	dst.Index(len).Set(reflect.ValueOf(e))

	return dst.Interface()
}
