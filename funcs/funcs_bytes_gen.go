package funcs

type (

	// All params is slice index, slice element
	// no return value
	MapFuncForBytes func(int, []byte)
	// return error to stop iterate
	MapErrFuncForBytes    func(int, []byte) error
	ApplyFuncForBytes     func(int, []byte) []byte
	ApplyErrFuncForBytes  func(int, []byte) ([]byte, error)
	FilterFuncForBytes    func(int, []byte) ([]byte, bool)
	FilterErrFuncForBytes ApplyErrFuncForBytes
)

// nil value for FuncT
var NILForBytes []byte = nil

// MapFor_FuncT iterate every elements of slice
func MapForBytes(slice [][]byte, fn MapFuncForBytes) {
	for index, o := range slice {
		fn(index, o)
	}
}

// MapWithErrFor_FuncT iterate every elements of slice, on error stop
func MapWithErrForBytes(slice [][]byte, fn MapErrFuncForBytes) (err error) {
	for index, o := range slice {
		if err = fn(index, o); err != nil {
			return
		}
	}
	return
}

// ApplyFor_FuncT apply function to every elements of slice
func ApplyForBytes(slice [][]byte, fn ApplyFuncForBytes) {
	for index, o := range slice {
		slice[index] = fn(index, o)
	}
}

// ApplyWithErrFor_FuncT apply function to every elements of slice, on error stop
func ApplyWithErrForBytes(slice [][]byte, fn ApplyErrFuncForBytes) (err error) {
	for index, o := range slice {
		if o, err = fn(index, o); err != nil {
			return
		}
		slice[index] = o
	}
	return
}

// FilterFor_FuncT iterate slice and filter with function, if function return true as useful
// append it to result slice
func FilterForBytes(slice [][]byte, fn FilterFuncForBytes) (res [][]byte) {
	MapForBytes(slice, func(index int, o []byte) {
		if o, use := fn(index, o); use {
			res = append(res, o)
		}
	})
	return res
}

// FilterWithErrFor_FuncT iterate slice and filter with function, if function return no error
// append it to result slice, else stop iterate
func FilterWithErrForBytes(slice [][]byte, fn FilterErrFuncForBytes) (res [][]byte, err error) {
	MapWithErrForBytes(slice, func(index int, o []byte) (e error) {
		if o, e = fn(index, o); err == nil {
			res = append(res, o)
		}
		return e
	})
	return
}

// ZipFor_FuncT zip two slice, if a slice is longer, the remains will match nil
func ZipForBytes(slice1, slice2 [][]byte) (res [][][]byte) {
	return zipForBytes(slice1, slice2, true)
}

// ZipShortFor_FuncT zip two slice, if a slice is longer, the remains will not be used
func ZipShortForBytes(slice1, slice2 [][]byte) (res [][][]byte) {
	return zipForBytes(slice1, slice2, false)
}

// zipFor_FuncT zip two slice, zipLong determin whether use the remains of longer slice
func zipForBytes(slice1, slice2 [][]byte, zipLong bool) (res [][][]byte) {
	var (
		i, l1, l2 int
		s1, s2    [][]byte
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
	res = make([][][]byte, 0, i)
	for i = 0; i < l2; i++ {
		res = append(res, [][]byte{s1[i], s2[i]})
	}
	if zipLong {
		for i = l2; i < l1; i++ {
			res = append(res, [][]byte{s1[i], NILForBytes})
		}
	}
	return res
}
