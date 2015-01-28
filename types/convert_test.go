package types

import (
	"bytes"
	"testing"

	"github.com/cosiner/golib/test"
)

func BenchmarkUnsafeBytesConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = UnsafeBytes("aaa")
	}
}

func BenchmarkNormalBytesConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = []byte("aaa")
	}
}

var bs = []byte("aaa")

func BenchmarkUnsafeStringConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = UnsafeString(bs)
	}
}

func BenchmarkNormalStringConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = string(bs)
	}
}

func TestUnsafeString(t *testing.T) {
	test.AssertEq(t, "UnsafeString", "abcde", UnsafeString([]byte("abcde")))
}

func TestUnsafeBytes(t *testing.T) {
	test.AssertEq(t, "UnsafeBytes", true, bytes.Equal([]byte("abcde"), UnsafeBytes("abcde")))
}

func TestHexStr2Uint(t *testing.T) {
	res, _ := HexStr2Uint("0xff")
	test.AssertEq(t, "HexStr2Uint1", uint(0xff), res)

	res, _ = HexStr2Uint("0x00fff")
	test.AssertEq(t, "HexStr2Uint2", uint(0xfff), res)
}

func TestReverseBits(t *testing.T) {
	tt := test.WrapTest(t)
	// 0000 0   0001 1  0010 2  0100 4  1000 8
	//          0011 3  0110 6  1100 C  1001 9
	tt.AssertEq("Reverse1", uint(0xcccccccc), ReverseBits(0x3333333300000000))
	tt.AssertEq("Reverse2", uint(0x48484848), ReverseBits(0x1212121200000000))
}

func TestReverseByte(t *testing.T) {
	tt := test.WrapTest(t)

	tt.AssertEq("ReverseByte1", uint8(0x12), ReverseByte(0x48))
	tt.AssertEq("ReverseByte2", uint8(0x3c), ReverseByte(0x3c))
}
