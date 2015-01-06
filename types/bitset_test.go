package types

import (
	"github.com/cosiner/golib/test"
	"testing"
)

func TestInOrNot(t *testing.T) {
	test.AssertEq(t, "In", uint(1<<2), In(2, uint((1<<0)|(1<<1)|(1<<2)|(1<<3)|(1<<4))))
	test.AssertEq(t, "NotIn", uint(1<<2), NotIn(2, uint((1<<0)|(1<<1)|(1<<3)|(1<<4))))
}
