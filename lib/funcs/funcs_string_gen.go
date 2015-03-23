// Package funcs implements a set of functions support for function programming
package funcs


// nil value for T
var NILForString string = ""

// MapFor_T iterate every elements of slice
func MapForString(slice []string, fn func(int, string)) {
	for index, o := range slice {
		fn(index, o)
	}
}

// MapWithErrFor_T iterate every elements of slice, on error stop
func MapWithErrForString(slice []string, fn func(int, string) error) (err error) {
	for index, o := range slice {
		if err = fn(index, o); err != nil {
			return
		}
	}
	return
}

// ApplyFor_T apply function to every elements of slice
func ApplyForString(slice []string, fn func(int, string) string) {
	for index, o := range slice {
		slice[index] = fn(index, o)
	}
}

// ApplyWithErrFor_T apply function to every elements of slice, on error stop
func ApplyWithErrForString(slice []string, fn func(int, string) (string, error)) (err error) {
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
func FilterForString(slice []string, fn func(int, string) (string, bool)) (res []string) {
	MapForString(slice, func(index int, o string) {
		if o, use := fn(index, o); use {
			res = append(res, o)
		}
	})
	return res
}

// FilterWithErrFor_T iterate slice and filter with function, if function return no error
// append it to result slice, else stop iterate
func FilterWithErrForString(slice []string, fn func(int, string) (string, error)) (res []string, err error) {
	MapWithErrForString(slice, func(index int, o string) (e error) {
		if o, e = fn(index, o); err == nil {
			res = append(res, o)
		}
		return e
	})
	return
}

// ZipFor_T zip two slice, if a slice is longer, the remains will match nil
func ZipForString(slice1, slice2 []string) (res [][]string) {
	return zipForString(slice1, slice2, true)
}

// ZipShortFor_T zip two slice, if a slice is longer, the remains will not be used
func ZipShortForString(slice1, slice2 []string) (res [][]string) {
	return zipForString(slice1, slice2, false)
}

// zipFor_T zip two slice, zipLong determin whether use the remains of longer slice
func zipForString(slice1, slice2 []string, zipLong bool) (res [][]string) {
	var (
		i, l1, l2 int
		s1, s2    []string
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
	res = make([][]string, 0, i)
	for i = 0; i < l2; i++ {
		res = append(res, []string{s1[i], s2[i]})
	}
	if zipLong {
		for i = l2; i < l1; i++ {
			res = append(res, []string{s1[i], NILForString})
		}
	}
	return res
}