package ast

import (
	"fmt"
	"testing"
)

// User is user
type User struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

func TestParseStruct(t *testing.T) {
	// tt := testing2.Wrap(t)
	c := Callback{
		Struct: func(a *Attrs) error {
			fmt.Println("Struct", a.TypeName, a.Field, a.Tag)
			return nil
		},
		Const: func(a *Attrs) error {
			fmt.Println("Const", a.TypeName, a.Name, a.Value)
			return nil
		},
		Func: func(a *Attrs) error {
			fmt.Println("Func", a.TypeName, a.Name)
			return nil
		},
	}
	ParseFile("ast.go", c)
}
