package types

import (
	. "github.com/cosiner/gohper/lib/generic"
)

// DelIndexObjFor_T delete all object that marked as true in delIndex
// from slice
func DelIndexObjFor_T(slice []T, delIndex []bool) []T {
	var newSlice []T
	for i, del := range delIndex {
		if !del {
			newSlice = append(newSlice, slice[i])
		}
	}
	return newSlice
}

// FilterV_T filter a slice, first return value is customed
func FilterV_T(filter func(int, T) (int, error), slice []T) (n int, err error) {
	var m int
	for index, s := range slice {
		if m, err = filter(index, s); err == nil {
			n += m
		} else {
			break
		}
	}
	return
}
