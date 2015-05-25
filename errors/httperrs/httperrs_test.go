package httperrs

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestErrors(t *testing.T) {
	tt := testing2.Wrap(t)

	tt.Nil(Server(nil))
	tt.Eq(500, ServerS("err").Code())
}
