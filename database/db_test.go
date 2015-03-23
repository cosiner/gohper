package database

import (
	"database/sql"
	"testing"

	"github.com/cosiner/golib/test"

	"github.com/cosiner/golib/types"
)

type User struct {
	Id   int
	Name string
}

func BenchmarkTypeInfo(b *testing.B) {
	db := NewDB()
	u := &User{}
	for i := 0; i < b.N; i++ {
		_ = db.TypeInfo(u)
	}
}

func BenchmarkFields(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = USER_ID | USER_NAME
	}
}

func BenchmarkFieldsCount(b *testing.B) {
	fs := USER_ID | USER_NAME
	for i := 0; i < b.N; i++ {
		_ = types.BitCountUint(fs)
	}
}

func BenchmarkSQLGet(b *testing.B) {
	db := NewDB()
	u := &User{}
	for i := 0; i < b.N; i++ {
		_ = db.TypeInfo(u).CacheGet(SELECT, USER_ID|USER_NAME, 0, SQLForSelect)
	}
}

func TestSQLCache(t *testing.T) {
	db := NewDB()
	u := &User{}
	sql := db.TypeInfo(u).CacheGet(SELECT, USER_ID|USER_NAME, 0, SQLForSelect)
	test.AssertEq(t, "SELECT id,name FROM user", sql)
}

func (u *User) NotFoundErr() error {
	return sql.ErrNoRows
}

func (u *User) DuplicateValueErr(key string) error {
	return nil
}

const (
	USER_ID uint = 1 << iota
	USER_NAME
	userFieldEnd
)

func (u *User) Table() string {
	return "user"
}

func (u *User) FieldValues(fields uint) []interface{} {
	vals, index := make([]interface{}, types.BitCountUint(fields)), 0
	if fields&USER_ID != 0 {
		vals[index] = u.Id
		index++
	}
	if fields&USER_NAME != 0 {
		vals[index] = u.Name
		index++
	}
	return vals[:index]
}

func (u *User) FieldPtrs(fields uint) []interface{} {
	vals, index := make([]interface{}, types.BitCountUint(fields)), 0
	if fields&USER_ID != 0 {
		vals[index] = &(u.Id)
		index++
	}
	if fields&USER_NAME != 0 {
		vals[index] = &(u.Name)
		index++
	}
	return vals[:index]
}

func (u *User) New() Model {
	return new(User)
}
