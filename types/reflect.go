package types

import (
	"reflect"
)

// IsSlice check whether or not param is slice
func IsSlice(s interface{}) bool {
	return s != nil && reflect.TypeOf(s).Kind() == reflect.Slice
}

// Equaler is a interface that compare whether two object is equal
type Equaler interface {
	EqualTo(interface{}) bool
}

// IndirectType return real type of value without pointer
func IndirectType(v interface{}) reflect.Type {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}
