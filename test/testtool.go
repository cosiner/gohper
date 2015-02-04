// Package test is a wrapper of testing that supply some useful functions for test
package test

import "testing"

// Test is a wrapper of testing.testing.TB
type Test struct {
	testing.TB
}

// WrapTest wrap testing.testing.TB or testing.B to Test
func WrapTest(t testing.TB) *Test {
	return &Test{t}
}

// AssertEq assert expect and got is equal, else print error message
func (t *Test) AssertEq(ident string, expect interface{}, got interface{}) {
	AssertEq(t.TB, ident, expect, got)
}

// AssertNE assert expect and got is not equal, else print error message
func (t *Test) AssertNE(ident string, expect interface{}, got interface{}) {
	AssertNE(t.TB, ident, expect, got)
}

// AssertTrue assert value is true
func (t *Test) AssertTrue(ident string, val bool) {
	AssertTrue(t.TB, ident, val)
}

// AssertFalse assert value is false
func (t *Test) AssertFalse(ident string, val bool) {
	AssertFalse(t.TB, ident, val)
}

// AssertEq assert expect and got is equal, else print error message
func AssertEq(t testing.TB, ident string, expect interface{}, got interface{}) {
	if expect != got {
		t.Errorf("Error in %s : expect %s, but got %s\n", ident, expect, got)
	}
}

// AssertNE assert expect and got is not equal, else print error message
func AssertNE(t testing.TB, ident string, expect interface{}, got interface{}) {
	if expect == got {
		t.Errorf("Error in %s : got same result: %s", ident, got)
	}
}

// AssertTrue assert value is true
func AssertTrue(t testing.TB, ident string, val bool) {
	AssertEq(t, ident, true, val)
}

// AssertFalse assert value is false
func AssertFalse(t testing.TB, ident string, val bool) {
	AssertEq(t, ident, false, val)
}
