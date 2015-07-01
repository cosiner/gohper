package httperrs

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestErrors(t *testing.T) {
	tt := testing2.Wrap(t)

	tt.Nil(Server.New(nil))
	tt.Nil(Server.NewS(""))

	err := Server.NewS("err")
	tt.Eq(err, Must(err))
	tt.Nil(Must(nil))

	tt.Eq(500, err.Code())

	_ = NewS("err", 400)
	_ = New(err, 400)
}
