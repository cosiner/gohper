package crypto

import (
	"testing"
)

func TestRandom(t *testing.T) {
	t.Log(RandAlphabet(32))
	t.Log(RandAlphanumeric(32))
	t.Log(RandAscii(32))
	t.Log(RandNumberal(32))
	t.Log(RandInt(1024))
	t.Log(RandInCharset(16, "0234d134"))
}
