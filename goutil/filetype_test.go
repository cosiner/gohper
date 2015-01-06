package goutil

import (
	"github.com/cosiner/golib/test"
	"testing"
)

func TestGoFileType(t *testing.T) {
	test.AssertEq(t, true, IsGoFile("aa.go"), "IsGoFile")
	test.AssertEq(t, true, IsGoSrcFile("aa.go"), "IsGoSrcFile")
	test.AssertEq(t, true, IsGoTestFile("aa_test.go"), "IsGoTestFile")
	test.AssertEq(t, false, IsGoSrcFile("aa_test.go"), "IsGoSrcFile")
	test.AssertEq(t, "aa_test.go", TestFileName("aa.go"), "TestFileName")
}
