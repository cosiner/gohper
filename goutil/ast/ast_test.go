package ast

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

type UserType int

const (
	ADMIN UserType = iota
	NORMAL
)

type UserIface interface {
	GetName() string
	SetName(string)
}

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
		Interface: func(a *Attrs) error {
			tt.Eq("UserIface", a.TypeName)
			tt.True(a.I.Method == "GetName" || a.I.Method == "SetName")
			return nil
		},
		Struct: func(a *Attrs) error {
			tt.Eq(a.TypeName, "User")
			if a.S.Field == "Name" {
				tt.Eq("string", a.S.Type)
				tt.Eq(`json:"name"`, string(a.S.Tag))
			} else if a.S.Field == "Id" {
				tt.Eq("int", a.S.Type)
				tt.Eq(`json:"id"`, string(a.S.Tag))
			} else {
				t.Fail()
			}
			return nil
		},

		Const: func(a *Attrs) error {
			tt.True(a.C.Name == "ADMIN" || a.C.Name == "NORMAL")
			tt.Eq(a.TypeName, "UserType")
			return nil
		},

		Func: func(a *Attrs) error {
			switch a.F.Name {
			case "GetName":
				tt.True(!a.F.PtrRecv)
				tt.Eq("User", a.TypeName)
			case "SetName":
				tt.True(a.F.PtrRecv)
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
