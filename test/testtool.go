// Package test is a wrapper of testing that supply some useful functions for test
package test

import (
	"testing"
)

// Test is a wrapper of testing.T
type Test struct {
	*testing.T
}

// WrapTest wrap testing.T to Test
func WrapTest(t *testing.T) *Test {
	return &Test{t}
}

// AssertEq assert expect and got is equal, else print error message
func (t *Test) AssertEq(ident string, expect interface{}, got interface{}) {
	AssertEq(t.T, ident, expect, got)
}

// AssertTrue assert value is true
func (t *Test) AssertTrue(ident string, val bool) {
	AssertTrue(t.T, ident, val)
}

// AssertEq assert expect and got is equal, else print error message
func AssertEq(t *testing.T, ident string, expect interface{}, got interface{}) {
	if expect != got {
		t.Errorf("Error in %s : expect %s, but got %s\n", ident, expect, got)
	}
}

// AssertTrue assert value is true
func AssertTrue(t *testing.T, ident string, val bool) {
	AssertEq(t, ident, true, val)
}
