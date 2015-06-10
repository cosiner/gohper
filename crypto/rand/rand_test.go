package rand

import (
	"testing"

	"github.com/cosiner/gohper/strings2"
	"github.com/cosiner/gohper/testing2"
	"github.com/cosiner/gohper/unsafe2"
)

func TestRandom(t *testing.T) {
	tt := testing2.Wrap(t)
	sfs := []func(int) (string, error){
		S.Alphabet, S.Alphanumeric, S.Numberal,
	}
	bfs := []func(int) ([]byte, error){
		B.Alphabet, B.Alphanumeric, B.Numberal,
	}
	cs := []string{
		ALPHABET, ALPHANUMERIC, NUMERALS,
	}

	for i := range sfs {
		testRandString(tt, sfs[i], cs[i])
		testRandBytes(tt, bfs[i], cs[i])
	}
}

func testRandString(tt testing2.TB, f func(int) (string, error), charset string) {
	s, e := f(32)
	tt.
		Eq(e, nil).
		Eq(32, len(s)).
		True(strings2.IsAllCharsIn(s, charset))
}

func testRandBytes(tt testing2.TB, f func(int) ([]byte, error), charset string) {
	s, e := f(32)
	tt.
		Eq(e, nil).
		Eq(32, len(s)).
		True(strings2.IsAllCharsIn(unsafe2.String(s), charset))
}

func BenchmarkRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		S.Alphabet(32)
	}
}
