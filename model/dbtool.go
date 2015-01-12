package model

import (
	"database/sql"
	"fmt"
)

// Insert insert model's field to database
func Insert(db *sql.DB, model Model, fields []Field, needId bool) (int64, error) {
	cols, ph := model.ColsSepPH(fields)
	sql := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", model.Table(), cols, ph)
	args := model.FieldVals(fields)
	return Exec(db, sql, args, true, needId)
}

// Update update model
func Update(db *sql.DB, model Model, fields []Field, whereFields []Field) (int64, error) {
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", model.Table(),
		model.ColsPH(fields), model.ColsPH(whereFields))
	args := model.FieldVals(append(fields, whereFields...))
	return Exec(db, sql, args, false, false)
}

// Delete delete model
func Delete(db *sql.DB, model Model, whereFields []Field) (int64, error) {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s", model.Table(), model.ColsPH(whereFields))
	args := model.FieldVals(whereFields)
	return Exec(db, sql, args, false, false)
}

// SelectOne select model from database
func SelectOne(db *sql.DB, model Model, fields []Field, whereField []Field) error {
	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s", model.Cols(fields), model.Table(),
		model.ColsPH(whereField))
	args := model.FieldVals(whereField)
	row := db.QueryRow(sql, args)
	return row.Scan(model.FieldPtrs(fields))
}

func Exec(db *sql.DB, sql string, args []interface{}, isInsert, needId bool) (int64, error) {
	res, err := db.Exec(sql, args)
	if err == nil {
		if isInsert {
			if needId {
				return res.LastInsertId()
			}
		} else {
			return res.RowsAffected()
		}
	}
	return 0, err
}
