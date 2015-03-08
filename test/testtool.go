// Package test is a wrapper of testing that supply some useful functions for test
package test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

// Test is a wrapper of testing.testing.TB
type Test struct {
	testing.TB
}

// WrapTest wrap testing.testing.TB or testing.B to Test
func WrapTest(t testing.TB) *Test {
	return &Test{t}
}

// AssertEq assert expect and got is equal, else print error message
func (t *Test) AssertEq(expect interface{}, got interface{}) {
	assertEq(t.TB, 1, expect, got)
}

// AssertNE assert expect and got is not equal, else print error message
func (t *Test) AssertNE(expect interface{}, got interface{}) {
	assertNE(t.TB, 1, expect, got)
}

// AssertTrue assert value is true
func (t *Test) AssertTrue(val bool) {
	assertEq(t.TB, 1, true, val)
}

// AssertFalse assert value is false
func (t *Test) AssertFalse(val bool) {
	assertEq(t.TB, 1, false, val)
}

// AssertNil assert value is nil
func (t *Test) AssertNil(value interface{}) {
	assertNil(t, 1, value)
}

// assertEq assert expect and got is equal, else print error message
func AssertEq(t testing.TB, expect interface{}, got interface{}) {
	assertEq(t, 1, expect, got)
}

// AssertNE assert expect and got is not equal, else print error message
func AssertNE(t testing.TB, expect interface{}, got interface{}) {
	assertNE(t, 1, expect, got)
}

// AssertNil assert value is nil
func AssertNil(t testing.TB, value interface{}) {
	assertNil(t, 1, value)
}

// AssertTrue assert value is true
func AssertTrue(t testing.TB, val bool) {
	assertEq(t, 1, true, val)
}

// AssertFalse assert value is false
func AssertFalse(t testing.TB, val bool) {
	assertEq(t, 1, false, val)
}

// assertEq assert expect and got is equal, else print error message
func assertEq(t testing.TB, skip int, expect interface{}, got interface{}) {
	if expect != got {
		t.Errorf("Error in %s : expect %s, but got %s\n",
			position(skip+1), fmt.Sprint(expect), fmt.Sprint(got))
	}
}

// assertNE assert expect and got is not equal, else print error message
func assertNE(t testing.TB, skip int, expect interface{}, got interface{}) {
	if expect == got {
		t.Errorf("Error in %s : expect different value, but got same: %s",
			position(skip+1), fmt.Sprint(got))
	}
}

// assertNil assert value is nil
func assertNil(t testing.TB, skip int, value interface{}) {
	if value != nil {
		t.Errorf("Error in %s: expect nil value, but got %s", position(skip+1), fmt.Sprint(value))
	}
}

func position(skip int) string {
	pc, file, line, _ := runtime.Caller(skip + 1)
	return filepath.Base(file) + ": " + filepath.Base(runtime.FuncForPC(pc).Name()) + ": " + strconv.Itoa(line)
}
