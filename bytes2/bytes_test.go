package bytes2

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestBytes(t *testing.T) {
	tt := testing2.Wrap(t)

	tt.DeepEq(TrimAfter([]byte("   ABCDE    # aaa"), []byte("#")), []byte("ABCDE"))
	tt.DeepEq(TrimBefore([]byte("   ABCDE    # aaa"), []byte("#")), []byte("aaa"))

	tt.DeepEq(SplitAndTrim([]byte(" A , B , C , D , E "), []byte(",")),
		[][]byte{[]byte("A"), []byte("B"), []byte("C"), []byte("D"), []byte("E")})

	tt.True(IsAllBytesIn([]byte("ABCDE"), []byte("ABCDEFG")))
	tt.False(IsAllBytesIn([]byte("ABCDEZ"), []byte("ABCDEFG")))

	tt.Eq(0, MakeBuffer(0, 8).Len())
}
