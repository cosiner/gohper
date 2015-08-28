package math

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestIndexDistance(t *testing.T) {
	tt := testing2.Wrap(t)
	index, indexUsed, remains := SegmentIndex([]int{2, 3}, 6)
	tt.Eq(-1, index)
	tt.Eq(0, indexUsed)
	tt.Eq(-1, remains)
}
