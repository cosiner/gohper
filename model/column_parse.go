package model

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cosiner/golib/types"
)

// colParser is a common column parser
type colParser struct {
	SqlBean
}

// Bind bind a sqlbean to parser
func (cp *colParser) Bind(sb SqlBean) {
	cp.SqlBean = sb
}

// MustValid check whether field is valid, otherwise panic
func (cp *colParser) MustValid(field Field) {
	if !cp.HasField(field) {
		cp.PanicUnknownField(field)
	}
}

// FieldCount return field count
func (cp *colParser) FieldCount() uint {
	return uint(len(cp.Fields()))
}

// ColValsExcp return column values exclude the excepts bitset
func (cp *colParser) FieldValsExcp(excepts []Field) []interface{} {
	return cp.FieldVals(cp.FieldsExcp(excepts))
}

// FieldsExcp return columns bitset exclude the Excp
func (cp *colParser) FieldsExcp(excepts []Field) []Field {
	var exists []Field
	var fs FieldSet = NewFieldSet(cp.FieldCount(), nil)
	for _, e := range excepts {
		cp.MustValid(e)
		fs.AddField(e)
	}
	for _, f := range cp.Fields() {
		if !fs.HasField(f) {
			exists = append(exists, f)
		}
	}
	return exists
}

// ColsAll return all columns joined with "," as string
func (cp *colParser) ColsAll() string {
	return cp.colsAll("", ",")
}

// ColsPHAll return all columns joined with ",",
// and each column append a placeholder suffix
func (cp *colParser) ColsPHAll() string {
	return cp.colsAll("=?", ",")
}

// colsAll return all columns string
func (cp *colParser) colsAll(suffix, sep string) string {
	colStr := strings.Join(cp.ColNames(), suffix+sep)
	if colStr == "" {
		return ""
	}
	return colStr + suffix
}

// Cols return columns string use given fieldset
func (cp *colParser) Cols(fields []Field) string {
	return cp.colsJoin("", ",", fields)
}

// ColsExcp return columns string exclude the excepts bitset
func (cp *colParser) ColsExcp(excepts []Field) string {
	return cp.Cols(cp.FieldsExcp(excepts))
}

// ColsSepPH return two string use given fieldset
// first string is columns, second string is placeholders
func (cp *colParser) ColsSepPH(fields []Field) (string, string) {
	fieldsStr := cp.colsJoin("", ",", fields)
	phStr := types.RepeatJoin("?", ",", len(fields))
	return fieldsStr, phStr
}

// ColsSepPHExcp return two string exclude given fieldset
func (cp *colParser) ColsSepPHExcp(excepts []Field) (string, string) {
	exist := cp.FieldsExcp(excepts)
	return cp.ColsSepPH(exist)
}

// ColsPH return columns string
// append each column with a placeholder '=?'
func (cp *colParser) ColsPH(fields []Field) string {
	return cp.colsJoin("=?", ",", fields)
}

// ColsPHExcp return columns string exclude the excepts bitset
// append each column with a placeholder '=?'
func (cp *colParser) ColsPHExcp(excepts []Field) string {
	return cp.ColsPH(cp.FieldsExcp(excepts))
}

// COLUMN_BUFSIZE if default buffer size to join columns
const COLUMN_BUFSIZE = 64

// colsJoin return column name exist in the exists bitset
// result like : col1+suffix+sep+col2+suffix+sep
func (cp *colParser) colsJoin(suffix, sep string, fields []Field) (col string) {
	if len(fields) != 0 {
		var buf *bytes.Buffer = bytes.NewBuffer(make([]byte, COLUMN_BUFSIZE))
		suffix = suffix + sep
		for _, f := range fields {
			cp.MustValid(f)
			buf.WriteString(cp.ColName(f))
			buf.WriteString(suffix)
		}
		colStr := buf.String()
		col = colStr[:len(colStr)-len(sep)]
	}
	return
}

// PanicUnknownField panic on unexpected field
func (cp *colParser) PanicUnknownField(field Field) {
	panic(fmt.Sprintf("Unexpected field %d for %s\n", field.UNum(), cp.Table()))
}

// // ColsPHVals return columns with placeholder, values
// func (cp *colParser) ColsPHVals(fields []Field) (col string,
// 	vals []interface{}) {
// 	return cp.colAndVals("=?", ",", fields)
// }

// // ColsSepPHVals return columns, placeholders, values
// func (cp *colParser) ColsSepPHVals(fields []Field) (
// 	col string, ph string, vals []interface{}) {
// 	col, vals = cp.colsAndVals("", ",", fields)
// 	ph = types.RepeatJoin("?", ",", len(fields))
// 	return
// }

// // ColsPHValsExcp return columns with placeholder, values
// func (cp *colParser) ColsPHValsExcp(excepts []Field) (col string,
// 	vals []interface{}) {
// 	return cp.colsAndValsExcp("=?", ",", excepts)
// }

// // ColsSepPHValsExcp return columns, placeholders, values
// func (cp *colParser) ColsSepPHValsExcp(excepts []Field) (col string,
// 	ph string, vals []interface{}) {
// 	col, vals = cp.colsAndValsExcp("", ",", excepts)
// 	ph = types.RepeatJoin("?", ",", cp.FieldCount-len(excepts))
// 	return
// }

// func (cp *colParser) colsAndVals(suffix, sep string, fields []Field) (
// 	col string, vals []interface{}) {
// 	buf := bytes.NewBuffer(make([]byte, COLUMN_BUFSIZE))
// 	vals := make([]interface{}, len(fields))
// 	suffix = suffix + sep
// 	for index, f := range fields {
// 		cp.MustValid(f)
// 		buf.WriteString(cp.ColName(f))
// 		buf.WriteString(suffix)
// 		vals[index] = cp.FieldVal(f)
// 	}
// 	if buf.Len() > 0 {
// 		colStr := buf.String()
// 		col = colStr[:len(colStr)-len(sep)]
// 	}
// 	return
// }

// func (cp *colParser) colsAndValsExcp(suffix, sep string, excepts []Field) (
// 	col string, vals []interface{}) {
// 	var fs FieldSet = NewFieldSet(cp.FieldCount())
// 	for _, e := range excepts {
// 		cp.MustValid(e)
// 		fs.AddField(e)
// 	}
// 	buf := bytes.NewBuffer(make([]byte, COLUMN_BUFSIZE))
// 	vals := make([]interface{}, len(fields))
// 	suffix = suffix + sep
// 	for _, f := range cp.Fields() {
// 		if !fs.HasField(f) {
// 			cp.MustValid(f)
// 			buf.WriteString(cp.ColName(f))
// 			buf.WriteString(suffix)
// 			vals[index] = cp.FieldVal(f)
// 		}
// 	}
// 	if buf.Len() > 0 {
// 		colStr := buf.String()
// 		col = colStr[:len(colStr)-len(sep)]
// 	}
// 	return
// }
