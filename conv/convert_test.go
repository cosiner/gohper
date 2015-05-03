package conv

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestHexStr2Uint(t *testing.T) {
	res, _ := Hex2Uint("0xff")
	testing2.Eq(t, uint64(0xff), res)

	res, _ = Hex2Uint("0x00fff")
	testing2.Eq(t, uint64(0xfff), res)
}

func TestReverseBits(t *testing.T) {
	tt := testing2.Wrap(t)
	// 0000 0   0001 1  0010 2  0100 4  1000 8
	//          0011 3  0110 6  1100 C  1001 9
	tt.Eq(uint(0xcccccccc), ReverseBits(0x3333333300000000))
	tt.Eq(uint(0x48484848), ReverseBits(0x1212121200000000))
}

func TestReverseByte(t *testing.T) {
	tt := testing2.Wrap(t)

	tt.Eq(uint8(0x12), ReverseByte(0x48))
	tt.Eq(uint8(0x3c), ReverseByte(0x3c))
}
