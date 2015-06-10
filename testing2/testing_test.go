package testing2

import (
	"strings"

	"testing"
)

func TestTest(t *testing.T) {
	tt := Wrap(t)
	var i []string
	// var j = []string{"1"}

	// Eq(t, 1, 1)
	// NE(t, t, nil)
	Nil(t, i)
	// NNil(t, j)
	// True(t, true)
	// False(t, false)
	// DeepEq(t, []string{"1"}, j)

	// tt.NNil("")
	// tt.NNil(1)
	// tt.NNil("a")
	// tt.NNil(struct{}{})
	// tt.Nil(nil)
	// tt.Eq(1, 1)
	// tt.NE(t, nil)
	tt.Nil(i)
	// tt.NNil(j)
	// tt.True(true)
	// tt.False(false)
	// tt.DeepEq([]string{"1"}, j)

	// defer Recover(t)

	// panic("panic")
}

func TestRecover(t *testing.T) {
	tt := Wrap(t)
	defer tt.Recover()
	panic("test")
}

func TestTestCase(t *testing.T) {
	Tests().
		Expect("abc").Arg("  abc   ").
		Expect("ab c").Arg("  ab c   ").
		Run(t, strings.TrimSpace)

}
