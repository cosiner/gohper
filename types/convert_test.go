package types

import (
	"bytes"
	"testing"

	"github.com/cosiner/golib/test"
)

func TestUnsafeString(t *testing.T) {
	test.AssertEq(t, "UnsafeString", "abcde", UnsafeString([]byte("abcde")))
}

func TestUnsafeBytes(t *testing.T) {
	test.AssertEq(t, "UnsafeBytes", true, bytes.Equal([]byte("abcde"), UnsafeBytes("abcde")))
}
