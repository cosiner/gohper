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
	colStr := strings.Join(cp.Columns(), suffix+sep)
	if colStr == "" {
		return ""
	}
	return colStr + suffix
}

// Cols return columns string use given fieldset
func (cp *colParser) Cols(fields ...Field) string {
	return cp.colsJoin("", ",", fields)
}

// ColsExcp return columns string exclude the Excps bitset
func (cp *colParser) ColsExcp(Excps ...Field) string {
	return cp.Cols(cp.colFieldsExcp(Excps)...)
}

// ColsSepPH return two string use given fieldset
// first string is columns, second string is placeholders
func (cp *colParser) ColsSepPH(fields ...Field) (string, string) {
	fieldsStr := cp.colsJoin("", ",", fields)
	phStr := types.RepeatJoin("?", ",", len(fields))
	return fieldsStr, phStr
}

// ColsSepPHExcp return two string exclude given fieldset
func (cp *colParser) ColsSepPHExcp(Excps ...Field) (string, string) {
	exist := cp.colFieldsExcp(Excps)
	return cp.ColsSepPH(exist...)
}

// ColsPH return columns string
// append each column with a placeholder '=?'
func (cp *colParser) ColsPH(fields ...Field) string {
	return cp.colsJoin("=?", ",", fields)
}

// ColsPHExcp return columns string exclude the Excps bitset
// append each column with a placeholder '=?'
func (cp *colParser) ColsPHExcp(Excps ...Field) string {
	return cp.ColsPH(cp.colFieldsExcp(Excps)...)
}

// ColVals return column values for given fields
func (cp *colParser) ColVals(fields ...Field) []interface{} {
	colVals := make([]interface{}, 0, len(fields))
	for _, f := range cp.Fields() {
		cp.MustValid(f)
		colVals = append(colVals, cp.FieldVal(f))
	}
	return colVals
}

// ColValsExcp return column values exclude the Excps bitset
func (cp *colParser) ColValsExcp(Excps ...Field) []interface{} {
	return cp.ColVals(cp.colFieldsExcp(Excps)...)
}

// colFieldsExcp return columns bitset exclude the Excp
func (cp *colParser) colFieldsExcp(Excps []Field) []Field {
	var exists []Field
	var fs FieldSet = NewFieldSet(cp.FieldCount())
	for _, e := range Excps {
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
			buf.WriteString(cp.ColumnName(f))
			buf.WriteString(suffix)
		}
		if buf.Len() != 0 {
			colStr := buf.String()
			col = colStr[:len(colStr)-len(sep)]
		}
	}
	return
}

// PanicUnknownField panic on unexpected field
func (cp *colParser) PanicUnknownField(field Field) {
	panic(fmt.Sprintf("Unexpected field %d for %s\n", field.UNum(), cp.Table()))
}
