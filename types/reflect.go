package types

import (
	"reflect"
)

// IsSlice check whether or not param is slice
func IsSlice(s interface{}) bool {
	return s != nil && reflect.TypeOf(s).Kind() == reflect.Slice
}
