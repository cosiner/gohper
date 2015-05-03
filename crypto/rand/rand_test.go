package rand

import (
	"testing"

	"github.com/cosiner/gohper/strings2"
	"github.com/cosiner/gohper/testing2"
)

func TestRandom(t *testing.T) {
	tt := testing2.Wrap(t)
	s, e := S.Alphabet(31)
	tt.True(e == nil)
	tt.Eq(31, len(s))
	tt.True(strings2.IsAllCharsIn(s, ALPHABET))
}

func BenchmarkRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		S.Alphabet(32)
	}
}
