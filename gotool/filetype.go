// Package gotool implements some utility functions related to go
package gotool

import (
	"mlib/util/types"
)

const (
	_SRC_SUFFIX  = ".go"
	_TEST_SUFFIX = "_test.go"
)

// IsGoSrcFile check whether file is go source file
func IsGoSrcFile(fname string) bool {
	return IsGoFile(fname) && !IsGoTestFile(fname)
}

// IsGoFile check whether file is go file
func IsGoFile(fname string) bool {
	return types.EndWith(fname, _SRC_SUFFIX)
}

// IsGoTestFile check whether file is go test file
func IsGoTestFile(fname string) bool {
	return types.EndWith(fname, _TEST_SUFFIX)
}

// CorrespondTestFile convert fname to corresponding test file name
func CorrespondTestFile(fname string) string {
	return fname[:len(fname)-len(_SRC_SUFFIX)] + _TEST_SUFFIX
}
