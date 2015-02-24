package crypto

import (
	"testing"

	"github.com/cosiner/golib/test"
)

var key = []byte("abcdefghijklmn")

func TestSecret(t *testing.T) {
	tt := test.WrapTest(t)
	enc := "abcde"
	val := SignSecret(key, enc)
	tt.Log(val)
	tt.AssertEq("1", enc, VerifySecret(key, val))
	key[0] = 'd'
	tt.AssertNE("2", enc, VerifySecret(key, val))
}
