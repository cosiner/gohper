// Package defval help reduce if-else block due that go has no '?:'
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
//
// It's not recommend to use if v is get from a function call.
func Nil(i, v interface{}) {
	ele := reflect.ValueOf(i).Elem()
	if ele.IsNil() {
		ele.Set(reflect.ValueOf(v))
	}
}

// NilFunc check whether *i is nil, if true set *i to f(),
// i must be a pointer, and f must be a function accept no parameters and return
// a value
//
// This function is not recommend to use.
func NilFunc(i, f interface{}) {
	ele := reflect.ValueOf(i).Elem()
	if ele.IsNil() {
		ele.Set(reflect.ValueOf(f).Call(nil)[0])
	}
}

// BoolStr check b, set *s to "true" or "false"
func BoolStr(b bool, s *string) {
	if b {
		*s = "true"
	} else {
		*s = "false"
	}
}

// BoolInt check b, set *i to 1 or 0
func BoolInt(b bool, i *int) {
	if b {
		*i = 1
	} else {
		*i = 0
	}
}
