package model

import (
	"database/sql"
	"errors"
)

const COLUMN_BUFSIZE = 64

var NotImplementedError = errors.New("Function not implemented.")

type SqlBean interface {
	// table name
	TableName() string
	// column names
	ColumnNames() []string
	// column values
	ColumnVals(fields uint) []interface{}
}

type SqlRunner interface {
	// add model, return sql and args
	Add(db *sql.DB) error
	// update model, return sql and args, if no field update all
	Update(db *sql.DB, fields int) error
	// delete model by id return sql and args
	Delete(db *sql.DB) error
	// select user return sql and args, if no field, by id
	Select(db *sql.DB, fields int) error
	// select limit
	SelectLimit(db *sql.DB, offset int, count int, fields int) error
}

type ColumnParser interface {
	Bind(sb SqlBean)
	// column name by index
	ColumnName(index int) string
	// ColumnCount
	ColumnCount() int
	// get column name by field
	Columns(fields uint) string
	//get column names with seperator
	ColumnsExcept(excepts uint) string
	// get column name with seperator and placeholder
	ColumnsPlaceHolder(fields uint) string
	// get column name with seperator and placeholder
	ColumnsPlaceHolderExcept(excepts uint) string
	// get column val by field name
	ColumnValsExcept(excepts uint) []interface{}
}

type Model interface {
	SqlBean
	ColumnParser
	SqlRunner
}

func NewColumnParser() ColumnParser {
	return &columnParse{}
}
