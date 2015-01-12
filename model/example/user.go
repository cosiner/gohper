package example

import (
	"database/sql"

	. "github.com/cosiner/golib/errors"
	. "github.com/cosiner/gomodule/model"
)

type Id uint32

type User struct {
	Id       Id
	Name     string
	Password string
	ColumnParser
}

func (u *User) Add(db *sql.DB) error {
	id, err := Insert(db, u, u.FieldsExcp(Fields(USER_ID)), true)
	if err == nil {
		u.Id = Id(id)
	}
	return err
}

func (u *User) Update(db *sql.DB, fields ...Field) error {
	Assert(u.Id != 0, Err("Update a non-exist user with id 0"))
	_, err := Update(db, u, fields, Fields(USER_ID))
	return err
}

func (u *User) Delete(db *sql.DB) error {
	Assert(u.Id != 0, Err("Delete a non-exist user with id 0"))
	_, err := Delete(db, u, Fields(USER_ID))
	return err
}

func (u *User) QueryById(db *sql.DB, fields ...Field) error {
	Assert(u.Id != 0, Err("Query a non-exist user with id 0"))
	return SelectOne(db, u, fields, Fields(USER_ID))
}
