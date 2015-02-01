package server

import (
	"strings"

	"github.com/cosiner/golib/test"

	"testing"
)

func TestRoute(t *testing.T) {
	tt := test.WrapTest(t)
	route := &route{
		Vars: []pathNode{
			pathNode{"user", 0},
			pathNode{"id", 1},
		},
	}
	path := []string{"abc", "123"}
	var user, id string
	route.ValuesScan(path, &user, &id)
	tt.AssertEq("1", "abc", user)
	tt.AssertEq("2", "123", id)
}

func TestCompile(t *testing.T) {
	tt := test.WrapTest(t)
	rt := new(router)
	path, vars, err := rt.compile("/:user/op/:id")
	if err != nil {
		panic(err)
	}
	tt.AssertTrue("3", err == nil)
	tt.AssertEq("4", "*/op/*", strings.Join(path, "/"))
	tt.AssertEq("5", 2, len(vars))
	tt.AssertEq("6", 0, vars[0].Index)
	tt.AssertEq("7", "user", vars[0].Var)
	tt.AssertEq("8", 2, vars[1].Index)
	tt.AssertEq("9", "id", vars[1].Var)
}
