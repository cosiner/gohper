package crypto

import (
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

var key = []byte("abcdefghijklmn")

func TestSecret(t *testing.T) {
	tt := test.WrapTest(t)
	enc := "abcde"
	val := SignSecret(key, enc)
	tt.Log(val)
	tt.AssertEq(enc, VerifySecret(key, val))
	key[0] = 'd'
	tt.AssertNE(enc, VerifySecret(key, val))
}
