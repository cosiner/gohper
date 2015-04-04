package example

import (
	"testing"

	"github.com/cosiner/gohper/database"

	"github.com/cosiner/gohper/lib/test"

	"github.com/cosiner/gohper/lib/types"
)

//go:generate gomodel -i $GOFILE
//go:generate gomodel -i $GOFILE -e
type User struct {
	Id   int
	Name string
}

type Name struct {
	First string
	Last  string
}

func BenchmarkTypeInfo(b *testing.B) {
	db := database.NewDB()
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
	db := database.NewDB()
	u := &User{}
	for i := 0; i < b.N; i++ {
		ti := db.TypeInfo(u)
		_ = ti.CacheGet(database.LIMIT_SELECT, USER_ID|USER_NAME, 0, ti.LimitSelectSQL)
	}
}

func TestSQLCache(t *testing.T) {
	db := database.NewDB()
	u := &User{}
	ti := db.TypeInfo(u)
	sql := ti.CacheGet(database.LIMIT_SELECT, USER_ID|USER_NAME, 0, ti.LimitSelectSQL)
	test.AssertEq(t, "SELECT id,name FROM user  LIMIT ?, ?", sql)
}

func TestCustom(t *testing.T) {
	db := database.NewDB()
	uti := db.TypeInfo(&User{})
	nti := db.TypeInfo(&Name{})
	t.Logf("SELECT %s,%s FROM %s,%s WHERE %s=%s\n",
		uti.TypedCols(USER_ID|USER_NAME),
		nti.TypedCols(NAME_FIRST|NAME_LAST),
		uti.Table,
		nti.Table,
		uti.TypedCols(USER_NAME), nti.TypedCols(NAME_FIRST))
}
