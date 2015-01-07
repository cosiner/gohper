package example

import "testing"

func TestAdd(t *testing.T) {
	user := (&User{1, "aaa", "ddd", "dddaa", "123", "1133222222", "dsaadqq", nil}).Init()
	t.Log("Columns:", user.Columns())
	t.Log("ColumnsStrAll:", user.ColumnsStrAll())
	t.Log("ColumnsStrPHExcept:", user.ColumnsPHStrExcept(user.FieldSet(USER_MOBILE, USER_ID)))
	t.Log("ColumnsStrExcept:", user.ColumnsStrExcept(user.FieldSet(USER_MOBILE, USER_ID)))
	t.Log("ColumnsStr:", user.ColumnsStr(user.FieldSet(USER_MOBILE, USER_ID)))
	t.Log("ColumnsStrPH:", user.ColumnsPHStr(user.FieldSet(USER_MOBILE, USER_ID)))
	t.Log("Table:", user.Table())
	t.Log("ColumnVals:", user.ColumnVals(user.FieldSet(USER_PASSWORD)))
	t.Log("ColumnValsExcept:", user.ColumnValsExcept(user.FieldSet(USER_PASSWORD)))
}
