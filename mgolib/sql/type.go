package mgolib

import (
	"reflect"
)

// whether or not struct field is unexported, according case of first character
func UnExported(f reflect.StructField) bool {
	return len(f.PkgPath) > 0
}

// whether or not a struct is embed in another struct
func Anonymous(f reflect.StructField) bool {
	return f.Anonymous
}

// whether it is struct kind
func StructKind(t reflect.Type) bool {
	return RealKind(t) == reflect.Struct
}

// RealKind
func RealKind(t reflect.Type) reflect.Kind {
	k := t.Kind()
	if k == reflect.Ptr {
		k = t.Elem().Kind()
	}
	return k
}

// Get RealType, if Ptr, return type point to, else itself
func RealType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
