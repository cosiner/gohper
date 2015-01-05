package test

import (
	"github.com/cosiner/gomodule/model"
	"testing"
)

func TestAdd(t *testing.T) {
	cp := model.NewColumnParser()
	user := &User{1, "aaa", "ddd", "dddaa", "123", "1133222222", "dsaadqq", cp}
	cp.Bind(user)
	t.Log(user.ColumnNames())
	t.Log(user.ColumnsExcept(USER_NONE))
	t.Log(user.ColumnsPlaceHolderExcept(USER_MOBILE | USER_PASSWORD))
	t.Log(user.TableName())
	t.Log(user.ColumnVals(USER_PASSWORD))
}
