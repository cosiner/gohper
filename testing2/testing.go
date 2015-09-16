// Package test is a wrapper of testing that supply some useful functions for test.
package testing2

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/cosiner/gohper/runtime2"
	"github.com/cosiner/gohper/terminal/ansi"
	"github.com/cosiner/gohper/terminal/color/output"
)

// TB is a wrapper of testing.testing.TB
type TB struct {
	testing.TB
}

// Wrap testing.testing.TB or testing.B to TB
func Wrap(t testing.TB) TB {
	return TB{t}
}

// Eq assert expect and got is equal, else print error message
func (t TB) Eq(expect, got interface{}) TB {
	eq(t.TB, 1, expect, got)

	return t
}

// DeepEq assert expect and got is deep-equal, else print error message
func (t TB) DeepEq(expect, got interface{}) TB {
	deepEq(t.TB, 1, expect, got)

	return t
}

// NE assert expect and got is not equal, else print error message
func (t TB) NE(expect, got interface{}) TB {
	ne(t.TB, 1, expect, got)

	return t
}

// True assert val is true
func (t TB) True(val bool) TB {
	eq(t.TB, 1, true, val)

	return t
}

// False assert val is false
func (t TB) False(val bool) TB {
	eq(t.TB, 1, false, val)

	return t
}

func (t TB) Nil(val interface{}) TB {
	nil_(t.TB, 1, val)

	return t
}

func (t TB) NNil(val interface{}) TB {
	nnil(t.TB, 1, val)

	return t
}

// Recover catch a panic and log it
func (t TB) Recover() TB {
	if e := recover(); e == nil {
		errorInfo(t.TB, 1, "panic", "not panic", false)
	}

	return t
}

// Recover catch a panic and log it
func (t TB) RecoverEq(s string) TB {
	if e := recover(); e == nil {
		errorInfo(t.TB, 1, "panic", "not panic", false)
	} else if es := fmt.Sprint(e); es != s {
		errorInfo(t, 1, "panic: "+s, "panic: "+es, false)
	}

	return t
}

// Eq assert expect and got is equal, else print error message
func Eq(t testing.TB, expect, got interface{}) TB {
	eq(t, 1, expect, got)

	return Wrap(t)
}

// DeepEq assert expect and got is deep-equal, else print error message
func DeepEq(t testing.TB, expect, got interface{}) TB {
	deepEq(t, 1, expect, got)

	return Wrap(t)
}

// NE assert expect and got is not equal, else print error message
func NE(t testing.TB, expect, got interface{}) TB {
	ne(t, 1, expect, got)

	return Wrap(t)
}

// True assert val is true
func True(t testing.TB, val bool) TB {
	eq(t, 1, true, val)

	return Wrap(t)
}

// False assert val is false
func False(t testing.TB, val bool) TB {
	eq(t, 1, false, val)

	return Wrap(t)
}

func Nil(t testing.TB, val interface{}) TB {
	nil_(t, 1, val)

	return Wrap(t)
}

func NNil(t testing.TB, val interface{}) TB {
	nnil(t, 1, val)

	return Wrap(t)
}

// eq assert expect and got is equal, else print error message
func eq(t testing.TB, skip int, expect, got interface{}) {
	if expect != got {
		errorInfo(t, skip+1, expect, got, true)
	}
}

// deepEq assert expect and got is deep-equal, else print error message
func deepEq(t testing.TB, skip int, expect, got interface{}) {
	if !reflect.DeepEqual(expect, got) {
		errorInfo(t, skip+1, expect, got, true)
	}
}

// ne assert expect and got is not equal, else print error message
func ne(t testing.TB, skip int, expect, got interface{}) {
	if expect == got {
		errorInfo(t, skip+1, "not equal", "equal", false)
	}
}

// nil_ assert val is nil
func nil_(t testing.TB, skip int, val interface{}) {
	if !isNil(val) {
		errorInfo(t, skip+1, "nil", fmt.Sprintf("%+v(%T)", val, val), false)
	}
}

// nnil assert val is not nil
func nnil(t testing.TB, skip int, val interface{}) {
	if isNil(val) {
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
	indexErrorInfo(t, skip+1, "", expect, got, withType)
}

func indexErrorInfo(t testing.TB, skip int, index string, expect, got interface{}, withType bool) {
	var (
		pos        = runtime2.Caller(skip + 1)
		exps, gots string
	)

	const formatT = "%+v(%T)"
	const format = "%+v"

	if !output.IsTTY {
		if withType {
			exps = fmt.Sprintf(formatT, expect, expect)
			gots = fmt.Sprintf(formatT, got, got)
		} else {
			exps = fmt.Sprintf(format, expect)
			gots = fmt.Sprintf(format, got)
		}
	} else {
		var red = ansi.Begin(ansi.Highlight, ansi.FgRed)
		var green = ansi.Begin(ansi.Highlight, ansi.FgGreen)
		var end = ansi.End()

		pos = ansi.Render(pos, ansi.Highlight, ansi.FgBlue)
		if withType {
			exps = fmt.Sprintf(green+formatT+end, expect, expect)
			gots = fmt.Sprintf(red+formatT+end, got, got)
		} else {
			exps = fmt.Sprintf(green+format+end, expect)
			gots = fmt.Sprintf(red+format+end, got)
		}
	}
	if index != "" {
		if output.IsTTY {
			index = ansi.Render(index, ansi.Highlight, ansi.FgYellow)
		}

		pos += ": " + index
	}

	t.Errorf("%s: expect: %s, got: %s", pos, exps, gots)
}

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}

	switch val := reflect.ValueOf(v); val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return val.IsNil()
	case reflect.Invalid:
		return fmt.Sprint(v) == "<nil>"
	default:
		return false
	}
}
