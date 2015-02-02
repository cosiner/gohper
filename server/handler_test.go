package server

import (
	"testing"

	"github.com/cosiner/golib/test"
)

func TestFuncHandler(t *testing.T) {
	tt := test.WrapTest(t)
	var fh Handler = new(funcHandler)
	_, is := fh.(MethodIndicator)
	tt.AssertTrue("H1", is)
}
