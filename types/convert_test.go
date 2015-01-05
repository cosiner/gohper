package types

import (
	"bytes"
	"mlib/util/test"
	"testing"
)

func TestUnsafeString(t *testing.T) {
	test.AssertEq(t, "abcde", UnsafeString([]byte("abcde")), "UnsafeString")
}

func TestUnsafeBytes(t *testing.T) {
	test.AssertEq(t, true, bytes.Equal([]byte("abcde"), UnsafeBytes("abcde")), "UnsafeBytes")
}
