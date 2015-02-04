// Package goutil implements some utility functions related to go
package goutil

import (
	"github.com/cosiner/golib/types"
)

const (
	srcSuffix  = ".go"
	testSuffix = "_test.go"
)

// IsGoSrcFile check whether file is go source file
func IsGoSrcFile(fname string) bool {
	return IsGoFile(fname) && !IsGoTestFile(fname)
}

// IsGoFile check whether file is go file
func IsGoFile(fname string) bool {
	return types.EndWith(fname, srcSuffix)
}

// IsGoTestFile check whether file is go test file
func IsGoTestFile(fname string) bool {
	return types.EndWith(fname, testSuffix)
}

// CorrespondTestFile convert fname to corresponding test file name
func CorrespondTestFile(fname string) string {
	return fname[:len(fname)-len(srcSuffix)] + testSuffix
}
