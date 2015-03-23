package crypto

import (
	"testing"
)

func TestRandom(t *testing.T) {
	t.Log(RandAlphabet(32))
	t.Log(RandAlphanumeric(32))
	t.Log(RandASCII(32))
	t.Log(RandNumberal(32))
	t.Log(RandInt(1024))
	t.Log(RandInCharset(16, "0234d134"))
}

func BenchmarkRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = RandAlphabet(32)
	}
}
