package test

import (
	"mlib/util/crypto"
	"testing"
)

func TestRandom(t *testing.T) {
	t.Log(crypto.RandAlphabet(32))
	t.Log(crypto.RandAlphanumeric(32))
	t.Log(crypto.RandAscii(32))
	t.Log(crypto.RandNumberal(32))
	t.Log(crypto.RandInt(1024))
	t.Log(crypto.RandInCharset(16, "0234d134"))
}
