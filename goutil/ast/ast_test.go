package ast

import (
	"testing"

	"github.com/cosiner/gohper/index"
	"github.com/cosiner/gohper/testing2"
)

type UserType int

const (
	ADMIN UserType = iota
	NORMAL
)

// User is user
type User struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

func (u User) GetName() string   { return "" }
func (u *User) SetName(n string) {}

func TestParse(t *testing.T) {
	tt := testing2.Wrap(t)

	c := Callback{
		Struct: func(a *Attrs) error {
			tt.Eq(a.TypeName, "User")
			return nil
		},
		StructField: func(a *Attrs) error {
			tt.True(index.StringIn(a.Field, []string{"Name", "Id"}) >= 0)
			tt.True(index.StringIn(a.Tag, []string{`json:"name"`, `json:"id"`}) >= 0)
			return nil
		},

		Const: func(a *Attrs) error {
			tt.True(a.Name == "ADMIN" || a.Name == "NORMAL")
			tt.Eq(a.TypeName, "UserType")
			return nil
		},

		Func: func(a *Attrs) error {
			switch a.Name {
			case "GetName":
				tt.True(!a.PtrRecv)
				tt.Eq("User", a.TypeName)
			case "SetName":
				tt.True(a.PtrRecv)
			case "TestParse":
				tt.Eq("", a.TypeName)
			default:
				tt.Fail()
			}
			return nil
		},
	}

	tt.True(ParseFile("ast_test.go", c) == nil)
}
