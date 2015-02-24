package database

import (
	"testing"

	"github.com/cosiner/golib/types"
)

type User struct {
	Id   int
	Name string
}

func BenchmarkTypeInfo(b *testing.B) {
	db := NewDB()
	for i := 0; i < b.N; i++ {
		_ = db.TypeInfo(&User{})
	}
}

func BenchmarkFields(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Fields(USER_ID, USER_NAME)
	}
}

func BenchmarkFieldsCount(b *testing.B) {
	fs := Fields(USER_ID, USER_NAME)
	for i := 0; i < b.N; i++ {
		_ = fs.BitCount()
	}
}

func TestSQLCache(t *testing.T) {
	db := NewDB()
	u := &User{}
	t.Log(db.CacheGet(db.SelectSQLCache, u, Fields(USER_ID, USER_NAME), EmptyFields, SQLForSelect))
}

const (
	USER_ID uint = iota
	USER_NAME
)

func (u *User) Table() string {
	return "user"
}

func (u *User) FieldValues(fields *types.LightBitSet) []interface{} {
	if fields == nil {
		return nil
	}
	vals, index := make([]interface{}, 0, fields.BitCount()), 0
	if fields.IsSet(USER_ID) {
		vals = append(vals, u.Id)
		index++
	}
	if fields.IsSet(USER_NAME) {
		vals = append(vals, u.Name)
		index++
	}
	return vals
}

func (u *User) FieldPtrs(fields *types.LightBitSet) []interface{} {
	if fields == nil {
		return nil
	}
	vals, index := make([]interface{}, 0, fields.BitCount()), 0
	if fields.IsSet(USER_ID) {
		vals = append(vals, &(u.Id))
		index++
	}
	if fields.IsSet(USER_NAME) {
		vals = append(vals, &(u.Name))
		index++
	}
	return vals
}

func (u *User) New() Model {
	return new(User)
}
