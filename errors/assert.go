package errors

// Assert val is true, else panic error
func Assert(val bool, err error) {
	if !val {
		panic(err)
	}
}
