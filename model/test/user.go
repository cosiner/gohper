package test

import (
	"github.com/cosiner/gomodule/model"
)

const (
	USER_TABLE = "user"
)
const (
	USER_ID uint = 1 << iota
	USER_NAME
	USER_PASSWORD
	USER_SALT
	USER_EMAIL
	USER_MOBILE
	USER_SHORTDESC
	USER_NONE uint = 0
)

var userColumns = []string{"id", "name", "password", "salt", "email", "mobile", "short_desc"}

type User struct {
	Id        int32
	Name      string
	Password  string
	Salt      string
	Email     string
	Mobile    string
	ShortDesc string
	model.ColumnParser
}

func (u *User) TableName() string {
	return USER_TABLE
}
func (b *User) ColumnNames() []string {
	return userColumns
}

func (u *User) ColumnVals(fields uint) []interface{} {
	colVals := make([]interface{}, 0, u.ColumnCount())
	if fields&USER_ID != 0 {
		colVals = append(colVals, u.Id)
	}
	if fields&USER_NAME != 0 {
		colVals = append(colVals, u.Name)
	}
	if fields&USER_EMAIL != 0 {
		colVals = append(colVals, u.Email)
	}
	if fields&USER_MOBILE != 0 {
		colVals = append(colVals, u.Mobile)
	}
	if fields&USER_PASSWORD != 0 {
		colVals = append(colVals, u.Password)
	}
	if fields&USER_SALT != 0 {
		colVals = append(colVals, u.Salt)
	}
	if fields&USER_SHORTDESC != 0 {
		colVals = append(colVals, u.ShortDesc)
	}
	if len(colVals) == 0 {
		return nil
	}
	return colVals
}
