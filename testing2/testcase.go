package testing2

import (
	"fmt"
	"reflect"
	"testing"
)

// NoCheck means don't check this value
var NoCheck = struct{}{}

// NonNil means this value should not be nil
var NonNil = struct{}{}

type TestCase struct {
	state   bool
	expects [][]interface{}
	args    [][]reflect.Value
}

func Tests() *TestCase {
	return &TestCase{
		state: wantExpect,
	}
}

func Expect(expect ...interface{}) *TestCase {
	return Tests().Expect(expect...)
}

const (
	wantArg    = true
	wantExpect = !wantArg
)

func (test *TestCase) checkState(state bool) {
	if test.state != state {
		panic("don't call same function continully")
	}

	test.state = !state
}

func (test *TestCase) Expect(expect ...interface{}) *TestCase {
	test.checkState(wantExpect)

	test.expects = append(test.expects, expect)
	return test
}

func (test *TestCase) True() *TestCase {
	return test.Expect(true)
}

func (test *TestCase) False() *TestCase {
	return test.Expect(false)
}

func (test *TestCase) Nil() *TestCase {
	return test.Expect(nil)
}

func (test *TestCase) Arg(args ...interface{}) *TestCase {
	test.checkState(wantArg)

	argVals := make([]reflect.Value, len(args))
	for i := range args {
		argVals[i] = reflect.ValueOf(args[i])
	}

	test.args = append(test.args, argVals)

	return test
}

func (test *TestCase) Run(t testing.TB, fn ...interface{}) *TestCase {
	for index := range test.expects {
		var expects = test.expects[index]

		var results = test.args[index]
		for _, f := range fn {
			results = reflect.ValueOf(f).Call(results)
		}

		for i, r := range results {
			if expect := expects[i]; expect != NoCheck {
				indexs := fmt.Sprintf("(%d:%d)", index+1, i+1)
				if expect == nil {
					indexNil(t, 1, indexs, "nil", r.Interface())
				} else if expect == NonNil {
					indexNil(t, 1, indexs, "not nil", "nil")
				} else {
					indexDeepEq(t, 1, indexs, expect, r.Interface())
				}
			}
		}
	}

	return test
}

func indexNil(t testing.TB, skip int, index string, expect, got interface{}) {
	if !isNil(got) {
		indexErrorInfo(t, skip+1, index, expect, got, false)
	}
}

func indexNNil(t testing.TB, skip int, index string, expect, got interface{}) {
	if isNil(got) {
		indexErrorInfo(t, skip+1, index, expect, got, false)
	}
}

// indexEq assert expect and got is equal, else print error message
func indexDeepEq(t testing.TB, skip int, index string, expect, got interface{}) {
	if !reflect.DeepEqual(expect, got) {
		indexErrorInfo(t, skip+1, index, expect, got, true)
	}
}
