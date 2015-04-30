// Package defval help for setup default value for primitive types
package defval

import "reflect"

func Int(v *int, d int) {
	if *v == 0 {
		*v = d
	}
}

func Int8(v *int8, d int8) {
	if *v == 0 {
		*v = d
	}
}

func Int16(v *int16, d int16) {
	if *v == 0 {
		*v = d
	}
}

func Int32(v *int32, d int32) {
	if *v == 0 {
		*v = d
	}
}

func Int64(v *int64, d int64) {
	if *v == 0 {
		*v = d
	}
}

func Uint(v *uint, d uint) {
	if *v == 0 {
		*v = d
	}
}

func Uint8(v *uint8, d uint8) {
	if *v == 0 {
		*v = d
	}
}

func Uint16(v *uint16, d uint16) {
	if *v == 0 {
		*v = d
	}
}

func Uint32(v *uint32, d uint32) {
	if *v == 0 {
		*v = d
	}
}

func Uint64(v *uint64, d uint64) {
	if *v == 0 {
		*v = d
	}
}

func String(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

// Nil check whether *i is nil, if true, set *i to v,
// i must be a pointer
func Nil(i, v interface{}) {
	ele := reflect.ValueOf(i).Elem()
	if ele.IsNil() {
		ele.Set(reflect.ValueOf(v))
	}
}
