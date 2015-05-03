package unsafe2

import (
	"bytes"
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func BenchmarkBytesConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Bytes("aaa")
	}
}

func BenchmarkNormalBytesConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = []byte("aaa")
	}
}

var bs = []byte("aaa")

func BenchmarkStringConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = String(bs)
	}
}

func BenchmarkNormalStringConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = string(bs)
	}
}

func TestString(t *testing.T) {
	testing2.Eq(t, "abcde", String([]byte("abcde")))
}

func TestBytes(t *testing.T) {
	testing2.Eq(t, true, bytes.Equal([]byte("abcde"), Bytes("abcde")))
}
