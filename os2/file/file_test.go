package file

import (
	"os"
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestCopyFile(t *testing.T) {
	tt := testing2.Wrap(t)
	src, dst := "/tmp/copytest", "/tmp/copytest.copy"
	if Create(src, nil) == nil {
		tt.True(Create(src, nil) != nil)
		tt.True(Copy(dst, src) == nil)
		tt.True(IsFile(dst))
		os.Remove(dst)
		os.Remove(src)
	}
}

func TestWriteFlag(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.Eq(os.O_APPEND, WriteFlag(false))
	tt.Eq(os.O_TRUNC, WriteFlag(true))
}
