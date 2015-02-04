package crypto

import (
	"bytes"
	"testing"

	"github.com/cosiner/golib/test"
)

var key = []byte("abcdefghijklmn")

func TestSecret(t *testing.T) {
	tt := test.WrapTest(t)
	enc := []byte("abcde")
	val := SignSecret(key, enc)
	tt.AssertTrue("1", bytes.Equal(enc, VerifySecret(key, val)))
	key[0] = 'd'
	tt.AssertFalse("1", bytes.Equal(enc, VerifySecret(key, val)))
}
