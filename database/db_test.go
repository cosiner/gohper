package database

import (
	"database/sql"
	"testing"

	"github.com/cosiner/gohper/lib/test"

	"github.com/cosiner/gohper/lib/types"
)

type User struct {
	Id   int
	Name string
}

type Name struct {
	First string
	Last  string
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
		ti := db.TypeInfo(u)
		_ = ti.CacheGet(LIMIT_SELECT, USER_ID|USER_NAME, 0, ti.LimitSelectSQL)
	}
}

func TestSQLCache(t *testing.T) {
	db := NewDB()
	u := &User{}
	ti := db.TypeInfo(u)
	sql := ti.CacheGet(LIMIT_SELECT, USER_ID|USER_NAME, 0, ti.LimitSelectSQL)
	test.AssertEq(t, "SELECT id,name FROM user  LIMIT ?, ?", sql)
}

func TestCustom(t *testing.T) {
	db := NewDB()
	uti := db.TypeInfo(&User{})
	nti := db.TypeInfo(&Name{})
	t.Logf("SELECT %s,%s FROM %s,%s WHERE %s\n",
		uti.TypedCols(USER_ID|USER_NAME),
		nti.TypedCols(NAME_FIRST|NAME_LAST),
		uti.Table,
		nti.Table,
		uti.TypedCol(USER_NAME)+"="+nti.TypedCol(NAME_FIRST))
}

const (
	USER_ID uint = 1 << iota
	USER_NAME
	userFieldEnd = iota
)

func (u *User) Table() string {
	return "user"
}

func (u *User) FieldValues(fields, reserve uint) []interface{} {
	vals, index := make([]interface{}, types.BitCountUint(fields)+int(reserve)), 0
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

func (u *User) NotFoundErr() error {
	return sql.ErrNoRows
}

func (u *User) DuplicateValueErr(key string) error {
	return nil
}

const (
	NAME_FIRST uint = 1 << iota
	NAME_LAST
	nameFieldEnd = iota
)

func (n *Name) Table() string {
	return "name"
}

func (n *Name) FieldValues(fields, reserve uint) []interface{} {
	vals, index := make([]interface{}, types.BitCountUint(fields)+int(reserve)), 0
	if fields&NAME_FIRST != 0 {
		vals[index] = n.First
		index++
	}
	if fields&NAME_LAST != 0 {
		vals[index] = n.Last
		index++
	}
	return vals[:index]
}

func (n *Name) FieldPtrs(fields uint) []interface{} {
	vals, index := make([]interface{}, types.BitCountUint(fields)), 0
	if fields&NAME_FIRST != 0 {
		vals[index] = &(n.First)
		index++
	}
	if fields&NAME_LAST != 0 {
		vals[index] = &(n.Last)
		index++
	}
	return vals[:index]
}

func (n *Name) New() Model {
	return new(Name)
}

func (n *Name) NotFoundErr() error {
	return sql.ErrNoRows
}

func (n *Name) DuplicateValueErr(key string) error {
	return nil
}
