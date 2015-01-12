// Package model defines a micro framework for database model
// for detail, see examples include in 'example/' to use with
// SqlBean and ColumnParser
// Because may  many code of sqlbean are most same, you can use
// my another tool 'github.com/cosiner/gotool/genmodel' to generate
// sqlbean's code, what it need is only a struct and it's field
package model

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
		FieldVals(field []Field) []interface{}
		// field's pointers
		FieldPtrs(field []Field) []interface{}
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
		// get column val by field name
		FieldValsExcp(excepts []Field) []interface{}
		// fields except given fields
		FieldsExcp(excepts []Field) []Field
		// all columns string : 'col1, col2, col3[]'
		ColsAll() string
		// all columns string with placeholder: 'col1=?, col2=?[]'
		ColsPHAll() string
		// get column name by field: 'col1, col2, col3[]'
		Cols(fields []Field) string
		//get column names with seperator: 'col1, col2, col3[]'
		ColsExcp(excepts []Field) string
		// get column name with seperator and placeholder: 'col1=?, col2=?[]'
		ColsPH(fields []Field) string
		// get column name with seperator and placeholder: 'col1=?, col2=?[]'
		ColsPHExcp(excepts []Field) string
		// fields string and placeholder string: ('col1, col2[]', '?, ?[]')
		ColsSepPH(fields []Field) (string, string)
		// fields string and placeholder string: ('col1, col2[]', '?, ?[]')
		ColsSepPHExcp(excepts []Field) (string, string)
		// // return columns string and values
		// ColsPHVals(fields []Field) (col string, vals []interface{})
		// // return columns, placeholders, values
		// ColsSepPHVals(fields []Field) (col string, ph string, vals []interface{})
		// // return columns string and values
		// ColsPHValsExcp(fields []Field) (col string, vals []interface{})
		// // return columns, placeholders, values
		// ColsSepPHValsExcp(fields []Field) (col string, ph string, vals []interface{})
	}

	Model interface {
		SqlBean
		ColumnParser
	}
)

// Fields recieve multi fields, return as slice
func Fields(fields ...Field) []Field {
	return fields
}

// NewColumnParser return the default column parser
func NewColumnParser() ColumnParser {
	return &colParser{}
}
