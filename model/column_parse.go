package model

import (
	"bytes"
	"mlib/util/types"
	"strings"
)

type columnParse struct {
	SqlBean
}

func (cp *columnParse) Bind(sb SqlBean) {
	cp.SqlBean = sb
}

func (cp *columnParse) ColumnCount() int {
	return len(cp.ColumnNames())
}

func (cp *columnParse) ColumnName(index int) (col string) {
	if count := cp.ColumnCount(); index >= 0 && count > index {
		cols := cp.ColumnNames()
		col = cols[index]
	}
	return
}

func (cp *columnParse) Columns(fields uint) string {
	return cp.columnsJoin(",", "", fields)
}

func (cp *columnParse) ColumnsExcept(excepts uint) string {
	return cp.Columns(cp.columnsIndexExcept(excepts))
}
func (cp *columnParse) ColumnsPlaceHolder(fields uint) string {
	return cp.columnsJoin(",", "=?", fields)
}
func (cp *columnParse) ColumnsPlaceHolderExcept(excepts uint) string {
	return cp.Columns(cp.columnsIndexExcept(excepts))
}

func (cp *columnParse) ColumnValsExcept(excepts uint) []interface{} {
	return cp.ColumnVals(cp.columnsIndexExcept(excepts))
}

func (cp *columnParse) ColumnsAll(sep string, suffix string) string {
	colStr := strings.Join(cp.ColumnNames(), suffix+sep)
	if colStr == "" {
		return ""
	}
	return colStr + suffix
}

// columnsIndexExcept return columns bitset exclude the except
func (cp *columnParse) columnsIndexExcept(excepts uint) uint {
	var exists uint
	for i := 0; i < cp.ColumnCount(); i++ {
		exists |= types.NotIn(i, excepts)
	}
	return exists
}

// ColumnsIn return column name exist in the exists bitset
// result like : col1+suffix+sep+col2+suffix+sep
func (cp *columnParse) columnsJoin(sep, suffix string, exists uint) string {
	if exists == 0 {
		return ""
	}
	var buf *bytes.Buffer
	suffix = suffix + sep
	for i, n := 0, cp.ColumnCount(); i < n; i++ {
		if types.In(i, exists) != 0 {
			if buf == nil {
				buf = bytes.NewBuffer(make([]byte, COLUMN_BUFSIZE))
			}
			buf.WriteString(cp.ColumnName(i))
			buf.WriteString(suffix)
		}
	}
	if buf == nil {
		return ""
	}
	colStr := buf.String()
	return colStr[:len(colStr)-len(sep)]
}
