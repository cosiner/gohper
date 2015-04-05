package crypto

import (
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

var key = []byte("abcdefghijklmn")

func TestSecret(t *testing.T) {
	tt := test.Wrap(t)
	enc := "abcde"
	val := SignSecret(key, enc)
	tt.Log(val)
	tt.Eq(enc, VerifySecret(key, val))
	key[0] = 'd'
	tt.NE(enc, VerifySecret(key, val))
}
