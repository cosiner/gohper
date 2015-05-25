package httperrs

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestErrors(t *testing.T) {
	tt := testing2.Wrap(t)

	tt.Nil(Server.New(nil))
	tt.Eq(503, Service.NewS("err").Code())
}
