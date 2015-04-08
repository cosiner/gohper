package database

import (
	"fmt"
	"reflect"
	"strings"

	ref "github.com/cosiner/gohper/lib/reflect"

	"github.com/cosiner/gohper/lib/goutil"
	"github.com/cosiner/gohper/lib/types"
)

type (
	// SQLCreator is type of sql creator
	SQLCreator func(fields, whereFields uint) string

	SQLType int

	// SQLCache is container of cached sql
	SQLCache map[uint]string

	Cacher []SQLCache

	// TypeInfo represent information of type
	// contains field count, table name, field names, field offsets
	TypeInfo struct {
		NumField uint
		Table    string
		Fields   []string
		Cacher
	}

	Cols interface {
		String() string
		Paramed() string
		OnlyParam() string
		Length() int
	}
	cols      []string
	singleCol string
)

const (
	INSERT SQLType = iota
	DELETE
	UPDATE
	LIMIT_SELECT
	defaultTypeEnd

	// _FIELD_SEP is seperator of columns
	_FIELD_SEP = ","
	// _FIELD_TAG is tag name of database column
	_FIELD_TAG    = "column"
	_FIELD_NOTCOL = "notcol"
)

var (
	printSQL = func(_ bool, _ string) {}
	// SQLTypeEnd used as all model's sql statement type count
	SQLTypeEnd      = defaultTypeEnd
	nilCols    Cols = make(cols, 0)
)

// SQLPrint enable sql print for each operation
func SQLPrint(enable bool, formatter func(formart string, v ...interface{})) {
	if enable {
		if formatter == nil {
			formatter = func(format string, v ...interface{}) {
				fmt.Printf(format, v...)
			}
		}
		printSQL = func(fromcache bool, sql string) {
			formatter("[SQL]CachedSQL:%t, sql:%s\n", fromcache, sql)
		}
	}
}

// FieldsExcp create fieldset except given fields
func FieldsExcp(numField uint, fields uint) uint {
	return (1<<numField - 1) & (^fields)
}

// FieldsIdentity create signature from fields
func FieldsIdentity(numField uint, fields, whereFields uint) uint {
	return fields<<numField | whereFields
}

// NewCacher create a common sql cacher with given type end,
// if type end is 0, use global type end
func NewCacher(typEnd SQLType) Cacher {
	if typEnd == 0 {
		typEnd = SQLTypeEnd
	}
	return make(Cacher, typEnd)
}

func (c Cacher) CacheGet(typ SQLType, id uint, create func() string) string {
	sql, has := c[typ][id]
	if has {
		printSQL(true, sql)
	} else {
		sql = create()
		c[typ][id] = sql
		printSQL(false, sql)
	}
	return sql
}

func (c Cacher) SQLTypeEnd(typ SQLType) Cacher {
	if int(typ) == len(c) {
		return c
	}
	cache := make([]SQLCache, typ)
	copy(cache, c)
	return cache
}

func (ti *TypeInfo) SQLTypeEnd(typ SQLType) {
	if int(typ) != len(ti.Cacher) {
		cache := make([]SQLCache, typ)
		copy(cache, ti.Cacher)
		ti.Cacher = cache
	}
}

// CacheGet get sql from cache container, if cache not exist, then create new
func (ti *TypeInfo) CacheGet(typ SQLType, fields, whereFields uint, create SQLCreator) (sql string) {
	cache := ti.Cacher[typ]
	id := FieldsIdentity(ti.NumField, fields, whereFields)
	if sql = cache[id]; sql == "" {
		sql = create(fields, whereFields)
		cache[id] = sql
		printSQL(false, sql)
	} else {
		printSQL(true, sql)
	}
	return
}

// InsertSQL create insert sql for given fields
func (ti *TypeInfo) InsertSQL(fields, _ uint) string {
	cols := ti.Cols(fields)
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)",
		ti.Table,
		cols,
		cols.OnlyParam())
}

// UpdateSQL create update sql for given fields
func (ti *TypeInfo) UpdateSQL(fields, whereFields uint) string {
	return fmt.Sprintf("UPDATE %s SET %s %s",
		ti.Table,
		ti.Cols(fields).Paramed(),
		ti.Where(whereFields))
}

// DeleteSQL create delete sql for given fields
func (ti *TypeInfo) DeleteSQL(_, whereFields uint) string {
	return fmt.Sprintf("DELETE FROM %s %s", ti.Table, ti.Where(whereFields))
}

// LimitSelectSQL create select sql for given fields
func (ti *TypeInfo) LimitSelectSQL(fields, whereFields uint) string {
	return fmt.Sprintf("SELECT %s FROM %s %s LIMIT ?, ?",
		ti.Cols(fields),
		ti.Table,
		ti.Where(whereFields))
}

// SQLForCount create select count sql
func (ti *TypeInfo) CountSQL(_, whereFields uint) string {
	return fmt.Sprintf("SELECT COUNT(*) FROM %s %s",
		ti.Table,
		ti.Where(whereFields))
}

func (ti *TypeInfo) Where(fields uint) string {
	if cols := ti.Cols(fields); cols.Length() != 0 {
		return "WHERE " + cols.Paramed()
	}
	return ""
}

func (ti *TypeInfo) TypedWhere(fields uint) string {
	if cols := ti.TypedCols(fields); cols.Length() != 0 {
		return "WHERE " + cols.Paramed()
	}
	return ""
}

// Cols return column names for given fields
// if fields is only one, return single column
// else return column slice
func (ti *TypeInfo) Cols(fields uint) Cols {
	return ti.colNames(fields, "")
}

// TypedCols return column names for given fields with type's table name as prefix
// like table.column
func (ti *TypeInfo) TypedCols(fields uint) Cols {
	return ti.colNames(fields, ti.Table+".")
}

func (ti *TypeInfo) colNames(fields uint, prefix string) Cols {
	fieldNames := ti.Fields
	if colCount := FieldCount(fields); colCount > 1 {
		names := make([]string, colCount)
		var index uint
		for i, l := uint(0), uint(len(fieldNames)); i < l; i++ {
			if (1<<i)&fields != 0 {
				names[index] = prefix + fieldNames[i]
				index++
			}
		}
		return cols(names[:index])
	} else if colCount == 1 {
		for i, l := uint(0), uint(len(fieldNames)); i < l; i++ {
			if (1<<i)&fields != 0 {
				return singleCol(prefix + fieldNames[i])
			}
		}
	}
	return nilCols
}

// it will first use field tag as column name, if no tag specified,
// use field name's camel_case
func parseTypeInfo(v Model) *TypeInfo {
	typ := ref.IndirectType(v)
	fieldNum := typ.NumField()
	fields := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		field := typ.Field(i)
		fieldName := field.Name
		if goutil.IsExported(fieldName) &&
			!strings.Contains(string(field.Tag), _FIELD_NOTCOL) &&
			!(field.Anonymous &&
				field.Type.Kind() == reflect.Struct) {
			if tagName := field.Tag.Get(_FIELD_TAG); tagName != "" {
				fieldName = tagName
			}
			fields = append(fields, types.SnakeString(fieldName))
		}
	}
	ti := &TypeInfo{
		NumField: uint(fieldNum),
		Table:    v.Table(),
		Fields:   fields,
		Cacher:   make(Cacher, SQLTypeEnd),
	}
	for i := SQLType(0); i < SQLTypeEnd; i++ {
		ti.Cacher[i] = make(SQLCache)
	}
	return ti
}

// String return columns string join with ",",
// result like "foo, bar"
func (c cols) String() string {
	return types.SuffixJoin(c, "", _FIELD_SEP)
}

// Paramed return columns string joind with "=?,", last "," was trimed,
// result like "foo=?, bar=?"
func (c cols) Paramed() string {
	return types.SuffixJoin(c, "=?", _FIELD_SEP)
}

// OnlyParam return columns placeholdered string, each column was replaced with "?"
// result like "?, ?, ?, ?", count of "?" is colums length
func (c cols) OnlyParam() string {
	return types.RepeatJoin("?", ",", len(c))
}

func (c cols) Length() int {
	return len(c)
}

// String return columns string join with ",",
// result like "foo, bar"
func (c singleCol) String() string {
	return string(c)
}

// Paramed return columns string joind with "=?,", last "," was trimed,
// result like "foo=?, bar=?"
func (c singleCol) Paramed() string {
	return string(c) + "=?"
}

// OnlyParam return columns placeholdered string, each column was replaced with "?"
// result like "?, ?, ?, ?", count of "?" is colums length
func (c singleCol) OnlyParam() string {
	return "?"
}

func (c singleCol) Length() int {
	return 1
}
