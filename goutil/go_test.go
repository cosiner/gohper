package goutil

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestFileType(t *testing.T) {
	testing2.Eq(t, true, IsGoFile("aa.go"))
	testing2.Eq(t, true, IsSrcFile("aa.go"))
	testing2.Eq(t, true, IsTestFile("aa_test.go"))
	testing2.Eq(t, false, IsSrcFile("aa_test.go"))
	testing2.Eq(t, "aa_test.go", SrcTestFile("aa.go"))
}
