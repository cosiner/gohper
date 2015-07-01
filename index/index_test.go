package index

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestIn(t *testing.T) {
	tt := testing2.Wrap(t)

	tt.Eq(1, StringIn("b", []string{"a", "b", "c"}))
	tt.Eq(0, StringIn("a", []string{"a", "b", "c"}))
	tt.Eq(-1, StringIn("d", []string{"a", "b", "c"}))

	tt.Eq(0, CharIn('a', "abc"))
	tt.Eq(1, CharIn('b', "abc"))
	tt.Eq(2, CharIn('d', "abd"))
	tt.Eq(-1, CharIn('e', "abcd"))

	testing2.
		Expect(uint(0)).Arg(-1, uint(1<<0|1<<2)).
		Expect(uint(0)).Arg(8, uint(1<<0|1<<2)).
		Expect(uint(0)).Arg(1, uint(1<<0|1<<2)).
		Expect(uint(1<<1)).Arg(1, uint(1<<1|1<<2|1<<3)).
		Expect(uint(1<<2)).Arg(2, uint(1<<2|1<<2|1<<3)).
		Run(t, BitIn)

	testing2.
		Expect(uint(0)).Arg(-1, uint(1<<0|1<<2)).
		Expect(uint(0)).Arg(2, uint(1<<0|1<<2)).
		Expect(uint(1<<9)).Arg(9, uint(1<<0|1<<2)).
		Expect(uint(1<<5)).Arg(5, uint(1<<1|1<<2|1<<3)).
		Expect(uint(1<<8)).Arg(8, uint(1<<2|1<<2|1<<3)).
		Run(t, BitNotIn)

	tt.Eq(2, ByteIn('A', 'B', 'C', 'A'))
	tt.Eq(0, ByteIn('A', 'A', 'B', 'C', 'A'))
	tt.Eq(-1, ByteIn('A', 'B', 'C'))

	tt.Eq(2, SortedNumberIn('C', 'A', 'B', 'C'))
	tt.Eq(0, SortedNumberIn('A', 'A', 'B', 'C', 'D'))
	tt.Eq(-1, SortedNumberIn('A', 'B', 'C'))

	tt.Eq(3, RuneIn('界', '你', '好', '世', '界'))
	tt.Eq(0, RuneIn('你', '你', '好', '世', '界'))
	tt.Eq(-1, RuneIn('是', '你', '好', '世', '界'))

}
