package reflect2

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestIncrSliceCap(t *testing.T) {
	tt := testing2.Wrap(t)
	sl := []string{"A"}

	sl = IncrAppend(sl, "B").([]string)
	tt.Eq(2, cap(sl))
	tt.DeepEq([]string{"A", "B"}, sl)
}
