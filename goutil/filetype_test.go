package goutil

import (
	"testing"

	"github.com/cosiner/golib/test"
)

func TestGoFileType(t *testing.T) {
	test.AssertEq(t, "IsGoFile", true, IsGoFile("aa.go"))
	test.AssertEq(t, "IsGoSrcFile", true, IsGoSrcFile("aa.go"))
	test.AssertEq(t, "IsGoTestFile", true, IsGoTestFile("aa_test.go"))
	test.AssertEq(t, "IsGoSrcFile", false, IsGoSrcFile("aa_test.go"))
	test.AssertEq(t, "TestFileName", "aa_test.go", CorrespondTestFile("aa.go"))
}
