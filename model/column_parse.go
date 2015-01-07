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

// FieldSet create a new fieldset
// the length is field count, panic on invalid field
func (cp *colParser) FieldSet(fields ...Field) FieldSet {
	fs := NewFieldSet(cp.FieldCount())
	for _, f := range fields {
		cp.MustEffectiveField(f)
		fs.AddField(f)
	}
	return fs
}

// MustEffectiveField check whether field is valid, otherwise panic
func (cp *colParser) MustEffectiveField(field Field) {
	if !cp.HasField(field) {
		cp.PanicUnknownField(field)
	}
}

// HasField check whether field is valid
// it's only iter the field list to check, if necessery, overwrite it
func (cp *colParser) HasField(field Field) bool {
	for _, f := range cp.Fields() {
		if f.Equal(field) {
			return true
		}
	}
	return false
}

// FieldCount return field count
func (cp *colParser) FieldCount() uint {
	return uint(len(cp.Fields()))
}

// ColumnsStrAll return all columns joined with "," as string
func (cp *colParser) ColumnsStrAll() string {
	return cp.columnsStrAll("", ",")
}

// ColumnsPHStrAll return all columns joined with ",",
// and each column append a placeholder suffix
func (cp *colParser) ColumnsPHStrAll() string {
	return cp.columnsStrAll("=?", ",")
}

// columnsStrAll return all columns string
func (cp *colParser) columnsStrAll(suffix, sep string) string {
	colStr := strings.Join(cp.Columns(), suffix+sep)
	if colStr == "" {
		return ""
	}
	return colStr + suffix
}

// ColumnsStr return columns string use given fieldset
func (cp *colParser) ColumnsStr(fields FieldSet) string {
	return cp.columnsJoin("", ",", fields)
}

// ColumnsStrExcept return columns string exclude the excepts bitset
func (cp *colParser) ColumnsStrExcept(excepts FieldSet) string {
	return cp.ColumnsStr(cp.columnFieldsExcept(excepts))
}

// ColumnsSepPHStr return two string use given fieldset
// first string is columns, second string is placeholders
func (cp *colParser) ColumnsSepPHStr(fields FieldSet) (string, string) {
	fieldsStr := cp.columnsJoin("", ",", fields)
	phStr := types.RepeatJoin("?", ",", fields.FieldCount())
	return fieldsStr, phStr
}

// ColumnsSepPHStrExcept return two string exclude given fieldset
func (cp *colParser) ColumnsSepPHStrExcept(excepts FieldSet) (string, string) {
	exist := cp.columnFieldsExcept(excepts)
	return cp.ColumnsSepPHStr(exist)
}

// ColumnsPHStr return columns string
// append each column with a placeholder '=?'
func (cp *colParser) ColumnsPHStr(fields FieldSet) string {
	return cp.columnsJoin("=?", ",", fields)
}

// ColumnsStrPHExcept return columns string exclude the excepts bitset
// append each column with a placeholder '=?'
func (cp *colParser) ColumnsPHStrExcept(excepts FieldSet) string {
	return cp.ColumnsPHStr(cp.columnFieldsExcept(excepts))
}

// ColumnVals return column values for given fields
func (cp *colParser) ColumnVals(fields FieldSet) []interface{} {
	colVals := make([]interface{}, 0, fields.BitCount())
	for _, f := range cp.Fields() {
		if fields.HasField(f) {
			colVals = append(colVals, cp.FieldVal(f))
		}
	}
	return colVals
}

// ColumnValsExcept return column values exclude the excepts bitset
func (cp *colParser) ColumnValsExcept(excepts FieldSet) []interface{} {
	return cp.ColumnVals(cp.columnFieldsExcept(excepts))
}

// columnFieldsExcept return columns bitset exclude the except
func (cp *colParser) columnFieldsExcept(excepts FieldSet) FieldSet {
	var exists FieldSet = NewFieldSet(cp.FieldCount())
	for _, f := range cp.Fields() {
		exists.ChangeField(f, !excepts.HasField(f))
	}
	return exists
}

// COLUMN_BUFSIZE if default buffer size to join columns
const COLUMN_BUFSIZE = 64

// columnsJoin return column name exist in the exists bitset
// result like : col1+suffix+sep+col2+suffix+sep
func (cp *colParser) columnsJoin(suffix, sep string, fields FieldSet) (col string) {
	if fields.Len() != 0 && fields.FieldCount() != 0 {
		var buf *bytes.Buffer = bytes.NewBuffer(make([]byte, COLUMN_BUFSIZE))
		suffix = suffix + sep
		for _, f := range cp.Fields() {
			if fields.HasField(f) {
				buf.WriteString(cp.ColumnName(f))
				buf.WriteString(suffix)
			}
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
