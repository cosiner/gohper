package example

import . "github.com/cosiner/gomodule/model"

var (
	userColumns = []string{"id", "name", "password", "salt", "email", "mobile", "short_desc"}
	userFields  = []Field{USER_ID, USER_NAME, USER_PASSWORD, USER_SALT, USER_EMAIL, USER_MOBILE, USER_SHORTDESC}
)

const (
	USER_ID Field = iota
	USER_NAME
	USER_PASSWORD
	USER_SALT
	USER_EMAIL
	USER_MOBILE
	USER_SHORTDESC
	userFieldEnd
	USER_TABLE = "user"
)

type User struct {
	Id        int32
	Name      string
	Password  string
	Salt      string
	Email     string
	Mobile    string
	ShortDesc string
	ColumnParser
}

func (u *User) Table() string {
	return USER_TABLE
}
func (u *User) Fields() []Field {
	return userFields
}

func (u *User) FieldVal(field Field) (val interface{}) {
	u.MustEffectiveField(field)
	switch field {
	case USER_ID:
		val = u.Id
	case USER_NAME:
		val = u.Name
	case USER_PASSWORD:
		val = u.Password
	case USER_SALT:
		val = u.Salt
	case USER_EMAIL:
		val = u.Email
	case USER_MOBILE:
		val = u.Mobile
	case USER_SHORTDESC:
		val = u.ShortDesc
	}
	return
}

func (u *User) Columns() []string {
	return userColumns
}

func (u *User) ColumnName(field Field) string {
	u.MustEffectiveField(field)
	return userColumns[field.Num()]
}

func (u *User) HasField(field Field) bool {
	return field < userFieldEnd
}

func (u *User) Init() *User {
	cp := NewColumnParser()
	cp.Bind(u)
	u.ColumnParser = cp
	return u
}
