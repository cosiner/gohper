package bytesize

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestSize(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.Eq(MustSize("1024"), uint64(1024))
	tt.Eq(MustSize("1024K"), uint64(1024*1024))
	tt.Eq(MustSize("1024M"), uint64(1024*1024*1024))
	tt.Eq(MustSize("1024G"), uint64(1024*1024*1024*1024))
}
