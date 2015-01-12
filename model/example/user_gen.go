package example

import (
	. "github.com/cosiner/gomodule/model"
)

const (
	USER_ID Field = iota
	USER_NAME
	USER_PASSWORD
	userFieldEnd
	USER_TABLE = "user"
)

var (
	userCols   = []string{"id", "name", "password"}
	userFields = []Field{USER_ID, USER_NAME, USER_PASSWORD}
)

// Init is
func (u *User) Init() *User {
	cp := NewColumnParser()
	cp.Bind(u)
	u.ColumnParser = cp
	return u
}

func (u *User) Table() string {
	return USER_TABLE
}

func (u *User) Fields() []Field {
	return userFields
}

func (u *User) HasField(field Field) bool {
	return field.UNum() < userFieldEnd.UNum()
}

func (u *User) ColNames() []string {
	return userCols
}

func (u *User) ColName(field Field) string {
	u.MustValid(field)
	return userCols[field.Num()]
}

func (u *User) FieldVals(fields []Field) (vals []interface{}) {
	fs := NewFieldSet(u.FieldCount(), u.MustValid, fields...)
	vals = make([]interface{}, len(fields))
	if fs.HasField(USER_ID) {
		vals = append(vals, u.Id)
	}
	if fs.HasField(USER_NAME) {
		vals = append(vals, u.Name)
	}
	if fs.HasField(USER_PASSWORD) {
		vals = append(vals, u.Password)
	}
	return
}

func (u *User) FieldPtrs(fields []Field) (ptrs []interface{}) {
	fs := NewFieldSet(u.FieldCount(), u.MustValid, fields...)
	ptrs = make([]interface{}, len(fields))
	if fs.HasField(USER_ID) {
		ptrs = append(ptrs, &(u.Id))
	}
	if fs.HasField(USER_NAME) {
		ptrs = append(ptrs, &(u.Name))
	}
	if fs.HasField(USER_PASSWORD) {
		ptrs = append(ptrs, &(u.Password))
	}
	return
}
