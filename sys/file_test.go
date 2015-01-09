package sys

import (
	"os"
	"testing"

	"github.com/cosiner/golib/test"
)

func TestCopyFile(t *testing.T) {
	tt := test.WrapTest(t)
	src, dst := "/tmp/copytest", "/tmp/copytest.copy"
	if CreateFor(src, nil) == nil {
		tt.AssertTrue("CopyFile", CreateFor(src, nil) != nil)
		tt.AssertTrue("CopyFile1", CopyFile(dst, src) == nil)
		tt.AssertTrue("CopyFile2", IsFile(dst))
		os.Remove(dst)
		os.Remove(src)
	}
}

func TestWriteFlag(t *testing.T) {
	tt := test.WrapTest(t)
	tt.AssertEq("WriteFlag", AP, WriteFlag(false))
	tt.AssertEq("WriteFlag", TC, WriteFlag(true))
}
