// Package test is a wrapper of testing that supply some useful functions for test.
package testing2

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/cosiner/gohper/runtime2"
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

// DeepEq assert expect and got is deep-equal, else print error message
func (t Test) DeepEq(expect interface{}, got interface{}) {
	deepEq(t.TB, 1, expect, got)
}

// NE assert expect and got is not equal, else print error message
func (t Test) NE(expect interface{}, got interface{}) {
	ne(t.TB, 1, expect, got)
}

// True assert val is true
func (t Test) True(val bool) {
	eq(t.TB, 1, true, val)
}

// False assert val is false
func (t Test) False(val bool) {
	eq(t.TB, 1, false, val)
}

func (t Test) Nil(val interface{}) {
	nil_(t.TB, 1, val)
}

func (t Test) NNil(val interface{}) {
	nnil(t.TB, 1, val)
}

// Recover catch a panic and log it
func (t Test) Recover() {
	if e := recover(); e == nil {
		errorInfo(t.TB, 1, "panic", "not panic", false)
	}
}

// Recover catch a panic and log it
func (t Test) RecoverEq(s string) {
	if e := recover(); e == nil {
		errorInfo(t.TB, 1, "panic", "not panic", false)
	} else if es := fmt.Sprint(e); es != s {
		errorInfo(t, 1, "panic: "+s, "panic: "+es, false)
	}
}

// Eq assert expect and got is equal, else print error message
func Eq(t testing.TB, expect interface{}, got interface{}) {
	eq(t, 1, expect, got)
}

// DeepEq assert expect and got is deep-equal, else print error message
func DeepEq(t testing.TB, expect interface{}, got interface{}) {
	deepEq(t, 1, expect, got)
}

// NE assert expect and got is not equal, else print error message
func NE(t testing.TB, expect interface{}, got interface{}) {
	ne(t, 1, expect, got)
}

// True assert val is true
func True(t testing.TB, val bool) {
	eq(t, 1, true, val)
}

// False assert val is false
func False(t testing.TB, val bool) {
	eq(t, 1, false, val)
}

func Nil(t testing.TB, val interface{}) {
	nil_(t, 1, val)
}

func NNil(t testing.TB, val interface{}) {
	nnil(t, 1, val)
}

// eq assert expect and got is equal, else print error message
func eq(t testing.TB, skip int, expect interface{}, got interface{}) {
	if expect != got {
		errorInfo(t, skip+1, expect, got, true)
	}
}

// deepEq assert expect and got is deep-equal, else print error message
func deepEq(t testing.TB, skip int, expect interface{}, got interface{}) {
	if !reflect.DeepEqual(expect, got) {
		errorInfo(t, skip+1, expect, got, true)
	}
}

// ne assert expect and got is not equal, else print error message
func ne(t testing.TB, skip int, expect interface{}, got interface{}) {
	if expect == got {
		errorInfo(t, skip+1, "not equal", "equal", false)
	}
}

// nil_ assert val is nil
func nil_(t testing.TB, skip int, val interface{}) {
	if !reflect.ValueOf(val).IsNil() {
		errorInfo(t, skip+1, "nil", "not nil", false)
	}
}

// nnil assert val is not nil
func nnil(t testing.TB, skip int, val interface{}) {
	if reflect.ValueOf(val).IsNil() {
		errorInfo(t, skip+1, "not nil", "nil", false)
	}
}

// Recover catch a panic and log it
func Recover(t testing.TB) {
	if e := recover(); e == nil {
		errorInfo(t, 1, "panic", "not panic", false)
	}
}

// Recover catch a panic and log it
func RecoverEq(t testing.TB, s string) {
	if e := recover(); e == nil {
		errorInfo(t, 1, "panic", "not panic", false)
	} else if es := fmt.Sprint(e); es != s {
		errorInfo(t, 1, "panic: "+s, "panic: "+es, false)
	}
}

func errorInfo(t testing.TB, skip int, expect, got interface{}, withType bool) {
	var (
		pos  = "\033[1;34m" + runtime2.Caller(skip+1) + "\033[0m"
		exps string
		gs   string
	)
	if withType {
		exps = fmt.Sprintf("\033[1;32m%+v(%T)\033[0m", expect, expect)
		gs = fmt.Sprintf("\033[1;31m%+v(%T)\033[0m", got, got)
	} else {
		exps = fmt.Sprintf("\033[1;32m%+v\033[0m", expect)
		gs = fmt.Sprintf("\033[1;31m%+v\033[0m", got)
	}

	t.Errorf("Error at %s : expect: %s, but got: %s", pos, exps, gs)
}
