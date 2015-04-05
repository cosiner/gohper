package sys

import (
	"os"
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

func TestCopyFile(t *testing.T) {
	tt := test.Wrap(t)
	src, dst := "/tmp/copytest", "/tmp/copytest.copy"
	if CreateFor(src, nil) == nil {
		tt.True(CreateFor(src, nil) != nil)
		tt.True(CopyFile(dst, src) == nil)
		tt.True(IsFile(dst))
		os.Remove(dst)
		os.Remove(src)
	}
}

func TestWriteFlag(t *testing.T) {
	tt := test.Wrap(t)
	tt.Eq(AP, WriteFlag(false))
	tt.Eq(TC, WriteFlag(true))
}
