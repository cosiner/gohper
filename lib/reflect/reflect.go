package reflect

import (
	"reflect"
	"strconv"

	"github.com/cosiner/gohper/lib/goutil"

	"github.com/cosiner/gohper/lib/errors"
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
	} else if !v.CanSet() {
		return errors.Err("Value can't be set")
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
		err = errors.Errorf("Unsupported type:%s", k.String())
	}
	return
}

// UnmarshalToStruct unmarshal map to struct, only primitive type will be unmarshaled
func UnmarshalToStruct(values map[string]string, v interface{}) error {
	value := reflect.ValueOf(v)
	kind := value.Kind()
	if kind != reflect.Ptr {
		return errors.Errorf("Non-pointer type: %s", kind.String())
	}
	value = value.Elem()
	kind = value.Kind()
	if kind != reflect.Struct {
		return errors.Errorf("Non-struct type:%s", kind.String())
	}
	for k, v := range values {
		if k == "" {
			continue
		}
		field := value.FieldByName(goutil.ExportedCase(k))
		if field.CanSet() {
			if err := UnmarshalPrimitive(v, field); err != nil {
				return err
			}
		}
	}
	return nil
}
