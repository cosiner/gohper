package mgolib

import (
	"testing"
)

func TestSql(t *testing.T) {
	qr := NewQueryRunner()
	qr.Select("Name", "Password").
		Select("Age").
		From("User", "Acc").
		Where(AND, EQ, "Id", 1).
		Where(OR, NOTNULL, "age").
		Orderby(Desc("id")).
		Groupby("aaa").
		Having(OR, EQ, "aaa", "1dddd").
		Where(AND, LIMIT, "", 2, 20).
		Where(AND, IN, "Age", 2, 3, 4)
	t.Log(qr.Sql())
}
