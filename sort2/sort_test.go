package sort2

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestBytes(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.DeepEq("abc", String("acb"))
	tt.DeepEq("abcd", string(Bytes([]byte("dcba"))))
}
