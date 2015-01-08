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
		// check whether has field
		HasField(field Field) bool
		// values of given field
		FieldVal(field Field) interface{}
		// column names
		ColNames() []string
		// column name by field
		ColName(field Field) string
	}
	ColumnParser interface {
		Bind(sb SqlBean)
		// check field effective, on error panic
		MustValid(field Field)
		// FieldCount
		FieldCount() uint
		// all columns string
		ColsAll() string
		// all columns string with placeholder
		ColsPHAll() string
		// get column name by field
		Cols(fields ...Field) string
		//get column names with seperator
		ColsExcp(excepts ...Field) string
		// get column name with seperator and placeholder
		ColsPH(fields ...Field) string
		// get column name with seperator and placeholder
		ColsPHExcp(excepts ...Field) string
		// fields string and placeholder string
		ColsSepPH(fields ...Field) (string, string)
		// fields string and placeholder string
		ColsSepPHExcp(excepts ...Field) (string, string)
		// ColumnVals return values of fields
		ColVals(fields ...Field) []interface{}
		// get column val by field name
		ColValsExcp(excepts ...Field) []interface{}
	}
	SqlRunner interface {
		// add model, return sql and args
		Add(db *sql.DB) error
		// update model, return sql and args, if no field update all
		Update(db *sql.DB, fields ...Field) error
		// delete model by id return sql and args
		Delete(db *sql.DB) error
		// select user return sql and args, if no field, by id
		Select(db *sql.DB, fields ...Field) error
		// select limit
		SelectLimit(db *sql.DB, offset int, count int, fields ...Field) error
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
