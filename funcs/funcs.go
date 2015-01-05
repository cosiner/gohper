package funcs

import (
	. "mlib/util/generic"
)

type (
	MapFuncFor_T func(int, T)
	// return error to stop iterate
	MapErrFuncFor_T    func(int, T) error
	ApplyFuncFor_T     func(int, T) T
	ApplyErrFuncFor_T  func(int, T) (T, error)
	FilterFuncFor_T    func(int, T) (T, bool)
	FilterErrFuncFor_T ApplyErrFuncFor_T
)

// nil value for T
var NILFor_T T = nil

// MapFor_T iterate every elements of slice
func MapFor_T(slice []T, fn MapFuncFor_T) {
	for index, o := range slice {
		fn(index, o)
	}
}

// MapWithErrFor_T iterate every elements of slice, on error stop
func MapWithErrFor_T(slice []T, fn MapErrFuncFor_T) (err error) {
	for index, o := range slice {
		if err = fn(index, o); err != nil {
			return
		}
	}
	return
}

// ApplyFor_T apply function to every elements of slice
func ApplyFor_T(slice []T, fn ApplyFuncFor_T) {
	for index, o := range slice {
		slice[index] = fn(index, o)
	}
}

// ApplyWithErrFor_T apply function to every elements of slice, on error stop
func ApplyWithErrFor_T(slice []T, fn ApplyErrFuncFor_T) (err error) {
	for index, o := range slice {
		if o, err = fn(index, o); err != nil {
			return
		}
		slice[index] = o
	}
	return
}

// FilterFor_T iterate slice and filter with function, if function return true as useful
// append it to result slice
func FilterFor_T(slice []T, fn FilterFuncFor_T) (res []T) {
	MapFor_T(slice, func(index int, o T) {
		if o, use := fn(index, o); use {
			res = append(res, o)
		}
	})
	return res
}

// FilterWithErrFor_T iterate slice and filter with function, if function return no error
// append it to result slice, else stop iterate
func FilterWithErrFor_T(slice []T, fn FilterErrFuncFor_T) (res []T, err error) {
	MapWithErrFor_T(slice, func(index int, o T) (e error) {
		if o, e = fn(index, o); err == nil {
			res = append(res, o)
		}
		return e
	})
	return
}

// ZipFor_T zip two slice, if a slice is longer, the remains will match nil
func ZipFor_T(slice1, slice2 []T) (res [][]T) {
	return zipFor_T(slice1, slice2, true)
}

// ZipShortFor_T zip two slice, if a slice is longer, the remains will not be used
func ZipShortFor_T(slice1, slice2 []T) (res [][]T) {
	return zipFor_T(slice1, slice2, false)
}

// zipFor_T zip two slice, zipLong determin whether use the remains of longer slice
func zipFor_T(slice1, slice2 []T, zipLong bool) (res [][]T) {
	var (
		i, l1, l2 int
		s1, s2    []T
	)
	l1, l2 = len(slice1), len(slice2)
	if l1 >= l2 {
		s1, s2 = slice1, slice2
	} else {
		l1, l2 = l2, l1
		s1, s2 = slice2, slice1
	}
	if l1 == 0 {
		return
	}
	i = l2
	if zipLong {
		i = l1
	}
	res = make([][]T, 0, i)
	for i = 0; i < l2; i++ {
		res = append(res, []T{s1[i], s2[i]})
	}
	if zipLong {
		for i = l2; i < l1; i++ {
			res = append(res, []T{s1[i], NILFor_T})
		}
	}
	return res
}
