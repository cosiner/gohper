package types

import (
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

func TestStr2Bytes(t *testing.T) {
	test.Eq(t, MustStr2Bytes("1024"), uint64(1024))
	test.Eq(t, MustStr2Bytes("1024K"), uint64(1024*1024))
	test.Eq(t, MustStr2Bytes("1024M"), uint64(1024*1024*1024))
	test.Eq(t, MustStr2Bytes("1024G"), uint64(1024*1024*1024*1024))
}
