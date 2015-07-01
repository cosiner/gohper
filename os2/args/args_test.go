package args

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestArgs(t *testing.T) {
	tt := testing2.Wrap(t)
	args := []string{
		"a", "1", "2", "3",
	}

	i, err := Int(args, 0, 0)
	tt.NNil(err)

	i, err = Int(args, 3, 3)
	tt.Nil(err)
	tt.Eq(3, i)

	i, err = Int(args, 4, 4)
	tt.Nil(err)
	tt.Eq(4, i)

	tt.Eq("a", String(args, 0, "0"))
	tt.Eq("4", String(args, 4, "4"))
}
