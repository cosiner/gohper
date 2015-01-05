package types

import (
	"github.com/cosiner/golib/test"
	"testing"
)

func TestStr2Bytes(t *testing.T) {
	test.AssertEq(t, MustStr2Bytes("1024"), 1024, "Muststr2Bytes")
	test.AssertEq(t, MustStr2Bytes("1024K"), 1024*1024, "Muststr2Bytes")
	test.AssertEq(t, MustStr2Bytes("1024M"), 1024*1024*1024, "Muststr2Bytes")
	test.AssertEq(t, MustStr2Bytes("1024G"), 1024*1024*1024*1024, "Muststr2Bytes")
}
