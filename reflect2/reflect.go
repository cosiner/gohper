package reflect2

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/cosiner/gohper/errors"
)

const (
	ErrNonPrimitive = errors.Err("not primitive type")
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

// UnmarshalPrimitive unmarshal bytes to primitive
func UnmarshalPrimitive(str string, v reflect.Value) (err error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch k := v.Kind(); k {
	case reflect.Bool:
		v.SetBool(str[0] == 't')
	case reflect.String:
		v.SetString(str)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if n, e := strconv.ParseInt(str, 10, 64); e == nil {
			v.SetInt(n)
		} else {
			err = e
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if n, e := strconv.ParseUint(str, 10, 64); e == nil {
			v.SetUint(n)
		} else {
			err = e
		}
	case reflect.Float32, reflect.Float64:
		if n, e := strconv.ParseFloat(str, v.Type().Bits()); e == nil {
			v.SetFloat(n)
		} else {
			err = e
		}
	default:
		return ErrNonPrimitive
	}
	return
}

func MarshalPrimitive(v reflect.Value) string {
	return fmt.Sprint(v.Interface())
}

func MarshalStruct(v interface{}, values map[string]string, tag string) {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		vfield := value.Field(i)
		if !vfield.CanInterface() {
			continue
		}
		tfield := typ.Field(i)
		name := tfield.Name
		if n := tfield.Tag.Get(tag); n == "-" {
			continue
		} else if n != "" {
			name = n
		} else {
			name = strings.ToLower(name)
		}
		values[name] = MarshalPrimitive(vfield)
	}
	return
}

func UnmarshalStruct(v interface{}, values map[string]string, tag string) {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	} else {
		panic("non-pointer type")
	}
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		vfield := value.Field(i)
		if !vfield.CanSet() {
			continue
		}
		tfield := typ.Field(i)
		name := tfield.Name
		if n := tfield.Tag.Get(tag); n == "-" {
			continue
		} else if n != "" {
			name = n
		} else {
			name = strings.ToLower(name)
		}
		UnmarshalPrimitive(values[name], vfield)
	}
	return
}
