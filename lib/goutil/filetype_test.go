package goutil

import (
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

func TestGoFileType(t *testing.T) {
	test.AssertEq(t, true, IsGoFile("aa.go"))
	test.AssertEq(t, true, IsGoSrcFile("aa.go"))
	test.AssertEq(t, true, IsGoTestFile("aa_test.go"))
	test.AssertEq(t, false, IsGoSrcFile("aa_test.go"))
	test.AssertEq(t, "aa_test.go", CorrespondTestFile("aa.go"))
}
