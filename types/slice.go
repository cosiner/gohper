package types

import (
	. "github.com/cosiner/golib/generic"
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
