package database

import (
	"testing"

	"github.com/cosiner/golib/test"
)

func TestParse(t *testing.T) {
	tt := test.WrapTest(t)
	type User struct {
		Id   int
		Name string
	}
	tt.Log(*Parse(User{}))
}
