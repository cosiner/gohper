package sys

import (
	"os"
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

func TestCopyFile(t *testing.T) {
	tt := test.WrapTest(t)
	src, dst := "/tmp/copytest", "/tmp/copytest.copy"
	if CreateFor(src, nil) == nil {
		tt.AssertTrue(CreateFor(src, nil) != nil)
		tt.AssertTrue(CopyFile(dst, src) == nil)
		tt.AssertTrue(IsFile(dst))
		os.Remove(dst)
		os.Remove(src)
	}
}

func TestWriteFlag(t *testing.T) {
	tt := test.WrapTest(t)
	tt.AssertEq(AP, WriteFlag(false))
	tt.AssertEq(TC, WriteFlag(true))
}
