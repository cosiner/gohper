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
