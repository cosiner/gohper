package bytesize

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestSize(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.Eq(uint64(1024), MustSize("1024"))
	tt.Eq(uint64(MB), MustSize("1024K"))
	tt.Eq(uint64(GB), MustSize("1024M"))
	tt.Eq(uint64(TB), MustSize("1024G"))
	tt.Eq(uint64(PB), MustSize("1024T"))
	tt.Eq(uint64(1024*PB), MustSize("1024P"))

	tt.Eq(uint64(0), SizeDef("-1", GB))
	tt.Eq(GB, SizeDef("abd", GB))

	defer tt.Recover()
	MustSize("abcde")
}
