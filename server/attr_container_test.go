package server

import (
	"testing"

	"github.com/cosiner/golib/test"
)

func TestUnlockedContainer(t *testing.T) {
	tt := test.WrapTest(t)
	c := NewAttrContainer()
	c.SetAttr("123", 231)
	tt.AssertEq("T1", 231, c.Attr("123"))
	tt.AssertFalse("T2", c.UpdateAttr("12", "a"))
	c.RemoveAttr("123")
	tt.AssertFalse("T3", c.UpdateAttr("123", "a"))
}

func TestLockedContainer(t *testing.T) {
	tt := test.WrapTest(t)
	c := NewLockedAttrContainer()
	c.SetAttr("123", 231)
	tt.AssertEq("T1", 231, c.Attr("123"))
	tt.AssertFalse("T2", c.UpdateAttr("12", "a"))
	c.RemoveAttr("123")
	tt.AssertFalse("T3", c.UpdateAttr("123", "a"))
}
