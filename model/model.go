package model

import "database/sql"

// CLOUMN_BUFSIZE is the default buffer size for column string

type (
	// Interface SqlBean define some functionss for a database object
	SqlBean interface {
		// table name
		Table() string
		// all fields
		Fields() []Field
		// values of given field
		FieldVal(field Field) interface{}
		// column names
		Columns() []string
		// column name by field
		ColumnName(field Field) string
	}
	ColumnParser interface {
		Bind(sb SqlBean)
		// create a new bitset use sqlbean's given fields
		FieldSet(fields ...Field) FieldSet
		// check whether has field
		HasField(field Field) bool
		// check field effective, on error panic
		MustEffectiveField(field Field)
		// FieldCount
		FieldCount() uint
		// all columns string
		ColumnsStrAll() string
		// all columns string with placeholder
		ColumnsPHStrAll() string
		// get column name by field
		ColumnsStr(fields FieldSet) string
		//get column names with seperator
		ColumnsStrExcept(excepts FieldSet) string
		// get column name with seperator and placeholder
		ColumnsPHStr(fields FieldSet) string
		// get column name with seperator and placeholder
		ColumnsPHStrExcept(excepts FieldSet) string
		// fields string and placeholder string
		ColumnsSepPHStr(fields FieldSet) (string, string)
		// fields string and placeholder string
		ColumnsSepPHStrExcept(excepts FieldSet) (string, string)
		// ColumnVals return values of fields
		ColumnVals(fields FieldSet) []interface{}
		// get column val by field name
		ColumnValsExcept(excepts FieldSet) []interface{}
	}
	SqlRunner interface {
		// add model, return sql and args
		Add(db *sql.DB) error
		// update model, return sql and args, if no field update all
		Update(db *sql.DB, fields FieldSet) error
		// delete model by id return sql and args
		Delete(db *sql.DB) error
		// select user return sql and args, if no field, by id
		Select(db *sql.DB, fields FieldSet) error
		// select limit
		SelectLimit(db *sql.DB, offset int, count int, fields FieldSet) error
	}
	Model interface {
		SqlBean
		ColumnParser
		SqlRunner
	}
)

func NewColumnParser() ColumnParser {
	return &colParser{}
}
