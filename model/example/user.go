package model

import (
	"database/sql"
	"fmt"

	. "github.com/cosiner/golib/errors"
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
	userCols   = [...]string{"id", "name", "password"}
	userFields = [...]Field{USER_ID, USER_NAME, USER_PASSWORD}
)

type User struct {
	Id       uint32
	Name     string
	Password string
}

// Init is
func (u *User) Init() *User {
	cp := NewColumnParser()
	cp.Bind(user)
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

func (u *User) FieldVal(field Field) (val interface{}) {
	u.MustValid(field)
	switch field {
	case USER_ID:
		val = u.Id
	case USER_NAME:
		val = u.Name
	case USER_PASSWORD:
		val = u.Password
	}
	return
}

func (u *User) Add(db *sql.DB) error {
	col, ph := u.ColsSepPHExcp(USER_ID)
	sql := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", u.Table(), col, ph)
	args := u.ColValsExcp(USER_ID)

	_, err := db.Exec(sql, args...)
	if err == nil {
		u.Id, _ = res.LastInsertId()
	}
	return err
}

func (u *User) Update(db *sql.DB, fields ...Field) error {
	Assert(u.Id != 0, Err("Update a non-exist user with id 0"))

	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s=?", u.Table(),
		u.ColsPH(fields), u.ColName(USER_ID))
	args := u.ColVals(fields...)

	_, err := db.Exec(sql, args...)
	return err
}

func (u *User) Delete(db *sql.DB) error {
	Assert(u.Id != 0, Err("Delete a non-exist user with id 0"))

	sql := fmt.Sprintf("DELETE FROM %s WHERE %s=?", u.Table(), u.ColName(USER_ID))
	args := u.ColumnVals(USER_ID)
	_, err := db.Exec(sql, args...)

	return err
}
