// Package test is a wrapper of testing that supply some useful functions for test
package test

import (
	"fmt"
	"testing"

	"github.com/cosiner/gohper/lib/runtime"
)

// Test is a wrapper of testing.testing.TB
type Test struct {
	testing.TB
}

// Wrap wrap testing.testing.TB or testing.B to Test
func Wrap(t testing.TB) Test {
	return Test{t}
}

// Eq assert expect and got is equal, else print error message
func (t Test) Eq(expect interface{}, got interface{}) {
	eq(t.TB, 1, expect, got)
}

// NE assert expect and got is not equal, else print error message
func (t Test) NE(expect interface{}, got interface{}) {
	ne(t.TB, 1, expect, got)
}

// True assert value is true
func (t Test) True(val bool) {
	eq(t.TB, 1, true, val)
}

// False assert value is false
func (t Test) False(val bool) {
	eq(t.TB, 1, false, val)
}

// Nil assert value is nil
func (t Test) Nil(value interface{}) {
	nil_(t, 1, value)
}

// NNil assert value is nil
func (t Test) NNil(value interface{}) {
	nnil(t, 1, value)
}

// Eq assert expect and got is equal, else print error message
func Eq(t testing.TB, expect interface{}, got interface{}) {
	eq(t, 1, expect, got)
}

// NE assert expect and got is not equal, else print error message
func NE(t testing.TB, expect interface{}, got interface{}) {
	ne(t, 1, expect, got)
}

// Nil assert value is nil
func Nil(t testing.TB, value interface{}) {
	nil_(t, 1, value)
}

func NNil(t testing.TB, value interface{}) {
	nnil(t, 1, value)
}

// True assert value is true
func True(t testing.TB, val bool) {
	eq(t, 1, true, val)
}

// False assert value is false
func False(t testing.TB, val bool) {
	eq(t, 1, false, val)
}

// eq assert expect and got is equal, else print error message
func eq(t testing.TB, skip int, expect interface{}, got interface{}) {
	if expect != got {
		t.Errorf("Error in %s : expect %s, but got %s\n",
			runtime.CallerPosition(skip+1), fmt.Sprint(expect), fmt.Sprint(got))
	}
}

// ne assert expect and got is not equal, else print error message
func ne(t testing.TB, skip int, expect interface{}, got interface{}) {
	if expect == got {
		t.Errorf("Error in %s : expect different value, but got same: %s",
			runtime.CallerPosition(skip+1), fmt.Sprint(got))
	}
}

// nil_ assert value is nil
func nil_(t testing.TB, skip int, value interface{}) {
	if value != nil {
		t.Errorf("Error in %s: expect nil value, but got %s", runtime.CallerPosition(skip+1), fmt.Sprint(value))
	}
}

// nnil assert value is not nil
func nnil(t testing.TB, skip int, value interface{}) {
	if value == nil {
		t.Errorf("Error in %s: expect non-nil value, but got nil", runtime.CallerPosition(skip+1))
	}
}
