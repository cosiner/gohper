package goutil

import (
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

func TestGoFileType(t *testing.T) {
	test.Eq(t, true, IsGoFile("aa.go"))
	test.Eq(t, true, IsGoSrcFile("aa.go"))
	test.Eq(t, true, IsGoTestFile("aa_test.go"))
	test.Eq(t, false, IsGoSrcFile("aa_test.go"))
	test.Eq(t, "aa_test.go", CorrespondTestFile("aa.go"))
}
