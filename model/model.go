// Package model defines a micro framework for database model
// for detail, see examples include in 'example/' to use with
// SqlBean and ColumnParser
// Because may  many code of sqlbean are most same, you can use
// my another tool 'github.com/cosiner/gotool/genmodel' to generate
// sqlbean's code, what it need is only a struct and it's field
package model

import "database/sql"

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
		// bind a sqlbean to parser
		Bind(sb SqlBean)
		// check field effective, on error panic
		// panic when sqlbean.HasField return false
		MustValid(field Field)
		// field count of sqlbean
		FieldCount() uint
		// all columns string : 'col1, col2, col3...'
		ColsAll() string
		// all columns string with placeholder: 'col1=?, col2=?...'
		ColsPHAll() string
		// get column name by field: 'col1, col2, col3...'
		Cols(fields ...Field) string
		//get column names with seperator: 'col1, col2, col3...'
		ColsExcp(excepts ...Field) string
		// get column name with seperator and placeholder: 'col1=?, col2=?...'
		ColsPH(fields ...Field) string
		// get column name with seperator and placeholder: 'col1=?, col2=?...'
		ColsPHExcp(excepts ...Field) string
		// fields string and placeholder string: ('col1, col2...', '?, ?...')
		ColsSepPH(fields ...Field) (string, string)
		// fields string and placeholder string: ('col1, col2...', '?, ?...')
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

// NewColumnParser return the default column parser
func NewColumnParser() ColumnParser {
	return &colParser{}
}
