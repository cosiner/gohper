package example

import (
	"testing"

	"github.com/cosiner/gohper/database"

	"github.com/cosiner/gohper/lib/test"

	"github.com/cosiner/gohper/lib/types"
)

//go:generate gomodel -i $GOFILE
type User struct {
	Id   int
	Name string
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
	db := database.New()
	u := &User{}
	for i := 0; i < b.N; i++ {
		ti := db.TypeInfo(u)
		_ = ti.CacheGet(database.LIMIT_SELECT, USER_ID|USER_NAME, 0, ti.LimitSelectSQL)
	}
}

func TestSQLCache(t *testing.T) {
	db := database.New()
	u := &User{}
	ti := db.TypeInfo(u)
	sql := ti.CacheGet(database.LIMIT_SELECT, USER_ID|USER_NAME, 0, ti.LimitSelectSQL)
	test.Eq(t, "SELECT id,name FROM user  LIMIT ?, ?", sql)
}

func TestCustom(t *testing.T) {
	db := database.New()
	uti := db.TypeInfo(&User{})
	nti := db.TypeInfo(&Name{})
	t.Logf("SELECT %s,%s FROM %s,%s WHERE %s=%s\n",
		uti.TypedCols(USER_ID|USER_NAME),
		nti.TypedCols(NAME_FIRST|NAME_LAST),
		uti.Table,
		nti.Table,
		uti.TypedCols(USER_NAME), nti.TypedCols(NAME_FIRST))
}







const (
    NAME_FIRST  uint = 1 << iota 
    NAME_LAST 
    nameFieldEnd = iota
)

func (n *Name) Table() string {
    return "name"
}

func (n *Name) FieldValues(fields,reserve uint) []interface{} {
    vals, index := make([]interface{}, types.BitCountUint(fields) + int(reserve)), 0
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

func (n *Name) New() database.Model {
    return new(Name)
}







const (
    USER_ID  uint = 1 << iota 
    USER_NAME 
    userFieldEnd = iota
)

func (u *User) Table() string {
    return "user"
}

func (u *User) FieldValues(fields,reserve uint) []interface{} {
    vals, index := make([]interface{}, types.BitCountUint(fields) + int(reserve)), 0
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

func (u *User) New() database.Model {
    return new(User)
}

