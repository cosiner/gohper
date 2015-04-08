package example

import (
	"testing"

	"github.com/cosiner/gohper/database"

	"github.com/cosiner/gohper/lib/test"
)

//go:generate gomodel -i $GOFILE
type User struct {
	Id  int
	Age int
}

type Name struct {
	First string
	Last  string
}

func BenchmarkTypeInfo(b *testing.B) {
	db := database.New()
	u := &User{}
	for i := 0; i < b.N; i++ {
		_ = db.TypeInfo(u)
	}
}

func BenchmarkSQLGet(b *testing.B) {
	db := database.New()
	u := &User{}
	for i := 0; i < b.N; i++ {
		ti := db.TypeInfo(u)
		_ = ti.CacheGet(database.LIMIT_SELECT, USER_ID|USER_AGE, 0, ti.LimitSelectSQL)
	}
}

func TestSQLCache(t *testing.T) {
	db := database.New()
	u := &User{}
	ti := db.TypeInfo(u)
	sql := ti.CacheGet(database.LIMIT_SELECT, USER_ID|USER_AGE, 0, ti.LimitSelectSQL)
	test.Eq(t, "SELECT id,age FROM user  LIMIT ?, ?", sql)
}

func TestCustom(t *testing.T) {
	db := database.New()
	uti := db.TypeInfo(&User{})
	nti := db.TypeInfo(&Name{})
	t.Logf("SELECT %s,%s FROM %s,%s WHERE %s=%s\n",
		uti.TypedCols(USER_ID|USER_ID),
		nti.TypedCols(NAME_FIRST|NAME_LAST),
		uti.Table,
		nti.Table,
		uti.TypedCols(USER_ID), nti.TypedCols(NAME_FIRST))
}

const (
	USER_ID uint = 1 << iota
	USER_AGE
	userFieldEnd = iota
)

func (u *User) Table() string {
	return "user"
}

func (u *User) Vals(fields uint, vals []interface{}) {
	if fields != 0 {
		index := 0
		if fields&USER_ID != 0 {
			vals[index] = u.Id
			index++
		}
		if fields&USER_AGE != 0 {
			vals[index] = u.Age
			index++
		}
	}
}

func (u *User) Ptrs(fields uint, ptrs []interface{}) {
	if fields != 0 {
		index := 0
		if fields&USER_ID != 0 {
			ptrs[index] = &(u.Id)
			index++
		}
		if fields&USER_AGE != 0 {
			ptrs[index] = &(u.Age)
			index++
		}
	}
}

func (u *User) New() database.Model {
	return new(User)
}

const (
	NAME_FIRST uint = 1 << iota
	NAME_LAST
	nameFieldEnd = iota
)

func (n *Name) Table() string {
	return "name"
}

func (n *Name) Vals(fields uint, vals []interface{}) {
	if fields != 0 {
		index := 0
		if fields&NAME_FIRST != 0 {
			vals[index] = n.First
			index++
		}
		if fields&NAME_LAST != 0 {
			vals[index] = n.Last
			index++
		}
	}
}

func (n *Name) Ptrs(fields uint, ptrs []interface{}) {
	if fields != 0 {
		index := 0
		if fields&NAME_FIRST != 0 {
			ptrs[index] = &(n.First)
			index++
		}
		if fields&NAME_LAST != 0 {
			ptrs[index] = &(n.Last)
			index++
		}
	}
}

func (n *Name) New() database.Model {
	return new(Name)
}
