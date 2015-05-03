// Package goutil implements some utility functions related to go
package goutil

import "strings"

const (
	srcSuffix  = ".go"
	testSuffix = "_test.go"
)

// IsSrcFile check whether file is go source file
func IsSrcFile(fname string) bool {
	return IsGoFile(fname) && !IsTestFile(fname)
}

// IsGoFile check whether file is go file
func IsGoFile(fname string) bool {
	return strings.HasSuffix(fname, srcSuffix)
}

// IsTestFile check whether file is go test file
func IsTestFile(fname string) bool {
	return strings.HasSuffix(fname, testSuffix)
}

// SrcTestFile convert fname to corresponding test file name
func SrcTestFile(fname string) string {
	return fname[:len(fname)-len(srcSuffix)] + testSuffix
}
