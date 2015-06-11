package testing2

import (
	"strings"
	"testing"
)

func TestTest(t *testing.T) {
	var i []string
	var j = []string{"1"}

	defer Wrap(t).
		Eq(1, 1).
		NE(t, nil).
		Nil(i).
		NNil(j).
		True(true).
		False(false).
		DeepEq([]string{"1"}, j).
		NNil("").
		NNil(1).
		NNil("a").
		NNil(struct{}{}).
		Nil(nil).
		Eq(1, 1).
		NE(t, nil).
		Nil(i).
		NNil(j).
		True(true).
		False(false).
		DeepEq([]string{"1"}, j).
		Recover()

	panic("panic")
}

func TestTestCase(t *testing.T) {
	Tests().
		Expect("abc").Arg("  abc   ").
		Expect("ab c").Arg("  ab c   ").
		Run(t, strings.TrimSpace)
}
