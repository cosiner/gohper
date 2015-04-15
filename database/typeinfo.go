package database

import (
	"database/sql"
	"fmt"
	"reflect"

	ref "github.com/cosiner/gohper/lib/reflect"

	"github.com/cosiner/gohper/lib/goutil"
	"github.com/cosiner/gohper/lib/types"
)

type (
	// SQLCreator create sql for model using fields and where fields
	SQLCreator func(fields, whereFields uint) string

	// TypeInfo represent information of type
	// contains field count, table name, field names, field offsets
	//
	// because count sql and limit select sql is both stored in the same cache,
	// you should not use a empty fields for limit select, that will conflict with
	// count sql and get the wrong sql statement.
	TypeInfo struct {
		NumField uint
		Table    string
		Fields   []string
		*Cacher
		prefix         string
		colsCache      map[uint]Cols
		typedColsCache map[uint]Cols
	}
)

const (
	// _FIELD_TAG is tag name of database column
	_FIELD_TAG = "column"
)

// FieldsExcp create fieldset except given fields
func FieldsExcp(numField uint, fields uint) uint {
	return (1<<numField - 1) & (^fields)
}

// FieldsIdentity create signature from fields
func FieldsIdentity(numField uint, fields, whereFields uint) uint {
	return fields<<numField | whereFields
}

// Stmt get sql from cache container, if cache not exist, then create new
func (ti *TypeInfo) Stmt(typ, fields, whereFields uint, create SQLCreator) (*sql.Stmt, error) {
	id := FieldsIdentity(ti.NumField, fields, whereFields)
	sql_, stmt := ti.Cacher.GetStmt(typ, id)
	if stmt == nil {
		sql_ = create(fields, whereFields)
		stmt, err := ti.Cacher.SetStmt(typ, id, sql_)
		if err == nil {
			printSQL(false, sql_)
		}
		return stmt, err
	}
	printSQL(true, sql_)
	return stmt, nil
}

func (ti *TypeInfo) InsertStmt(fields uint) (*sql.Stmt, error) {
	return ti.Stmt(LIMIT_SELECT, fields, 0, ti.InsertSQL)
}

func (ti *TypeInfo) UpdateStmt(fields, whereFields uint) (*sql.Stmt, error) {
	return ti.Stmt(UPDATE, fields, whereFields, ti.UpdateSQL)
}

func (ti *TypeInfo) DeleteStmt(whereFields uint) (*sql.Stmt, error) {
	return ti.Stmt(DELETE, 0, whereFields, ti.DeleteSQL)
}

func (ti *TypeInfo) LimitSelectStmt(fields, whereFields uint) (*sql.Stmt, error) {
	return ti.Stmt(LIMIT_SELECT, fields, whereFields, ti.LimitSelectSQL)
}

func (ti *TypeInfo) SelectOneStmt(fields, whereFields uint) (*sql.Stmt, error) {
	return ti.Stmt(SELECT_ONE, fields, whereFields, ti.SelectOneSQL)
}

func (ti *TypeInfo) CountStmt(whereFields uint) (*sql.Stmt, error) {
	return ti.Stmt(LIMIT_SELECT, 0, whereFields, ti.CountSQL)
}

// InsertSQL create insert sql for given fields
func (ti *TypeInfo) InsertSQL(fields, _ uint) string {
	cols := ti.Cols(fields)
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)",
		ti.Table,
		cols.String(),
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

// LimitSelectSQL create select sql for given fields
func (ti *TypeInfo) SelectOneSQL(fields, whereFields uint) string {
	return fmt.Sprintf("SELECT %s FROM %s %s LIMIT 1",
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
	cols := ti.colsCache[fields]
	if cols == nil {
		cols = ti.colNames(fields, "")
		ti.colsCache[fields] = cols
	}
	return cols
}

// TypedCols return column names for given fields with type's table name as prefix
// like table.column
func (ti *TypeInfo) TypedCols(fields uint) Cols {
	cols := ti.typedColsCache[fields]
	if cols == nil {
		cols = ti.colNames(fields, ti.prefix)
		ti.typedColsCache[fields] = cols
	}
	return cols
}

// colNames get fields names, each field prepend a prefix string
// if fields count is 0, type nilCols was returned
// if fields count is 1, type singleCols was returned
// otherwise, type cols was returned
func (ti *TypeInfo) colNames(fields uint, prefix string) Cols {
	fieldNames := ti.Fields
	if colCount := FieldCount(fields); colCount > 1 {
		names := make([]string, colCount)
		var index int
		for i, l := uint(0), uint(len(fieldNames)); i < l; i++ {
			if (1<<i)&fields != 0 {
				names[index] = ti.Table + "." + fieldNames[i]
				index++
			}
		}
		return &cols{cols: names}
	} else if colCount == 1 {
		for i, l := uint(0), uint(len(fieldNames)); i < l; i++ {
			if (1<<i)&fields != 0 {
				return singleCol(prefix + fieldNames[i])
			}
		}
	}
	return zeroCols
}

// parseTypeInfo will first use field tag as column name, the tag key is 'column',
// if no tag specified, use field name's camel_case, disable a field by put 'notcol'
// in field tag
func parseTypeInfo(v Model, db *DB) *TypeInfo {
	typ := ref.IndirectType(v)
	fieldNum := typ.NumField()
	fields := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		field := typ.Field(i)
		fieldName := field.Name
		// Exported + !(anonymous && structure)
		if goutil.IsExported(fieldName) &&
			!(field.Anonymous &&
				field.Type.Kind() == reflect.Struct) {

			tagName := field.Tag.Get(_FIELD_TAG)
			if tagName == "-" {
				continue
			}
			if tagName != "" {
				fieldName = tagName
			}
			fields = append(fields, types.SnakeString(fieldName))
		}
	}
	return &TypeInfo{
		NumField:       uint(fieldNum),
		Table:          v.Table(),
		Fields:         fields,
		Cacher:         NewCacher(SQLTypes, db),
		prefix:         v.Table() + ".",
		colsCache:      make(map[uint]Cols),
		typedColsCache: make(map[uint]Cols),
	}
}
