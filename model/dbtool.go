package model

import (
	"database/sql"
	"fmt"
)

var printSql = func(_ string) {}

// EnableSqlPrint enable print sql
func EnableSqlPrint() {
	printSql = func(sql string) { fmt.Println(sql) }
}

// Insert insert model's field to database
func Insert(db *sql.DB, model Model, fields []Field, needId bool) (int64, error) {
	cols, ph := model.ColsSepPH(fields)
	sql := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", model.Table(), cols, ph)
	args := model.FieldVals(fields)
	return Exec(db, sql, args, needId)
}

// Update update model
func Update(db *sql.DB, model Model, fields []Field, whereFields []Field) (int64, error) {
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", model.Table(),
		model.ColsPH(fields), model.ColsPH(whereFields))
	args := model.FieldVals(append(fields, whereFields...))
	return Exec(db, sql, args, false)
}

// Delete delete model
func Delete(db *sql.DB, model Model, whereFields []Field) (int64, error) {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s", model.Table(), model.ColsPH(whereFields))
	args := model.FieldVals(whereFields)
	return Exec(db, sql, args, false)
}

// SelectOne select model from database
func SelectOne(db *sql.DB, model Model, fields []Field, whereField []Field) error {
	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s", model.Cols(fields), model.Table(),
		model.ColsPH(whereField))
	args := model.FieldVals(whereField)
	row := db.QueryRow(sql, args...)
	printSql(sql)
	return row.Scan(model.FieldPtrs(fields)...)
}

// ExecInsert perform insert sql, if needId return last insert id,
// else affected row count
func ExecInsert(db *sql.DB, s string, args []interface{}, needId bool) (int64,
	error) {
	return Exec(db, s, args, needId)
}

// ExecUpdate perform update or delete sql, return affected row count
func ExecUpdate(db *sql.DB, s string, args []interface{}) (int64, error) {
	return Exec(db, s, args, false)
}

// Exec execute a sql
func Exec(db *sql.DB, s string, args []interface{}, needId bool) (ret int64, err error) {
	printSql(s)
	res, err := db.Exec(s, args...)
	if err == nil {
		ret, err = ResolveResult(res, needId)
	}
	return
}

// BatchExec bached execute sql, the parameter function generate arguments
func BatchExec(db *sql.DB, s string, fn func(error) []interface{}, needId, needTx bool) (ret []int64, err error) {
	var res sql.Result
	stmt, err, deferFn := TxOrNot(db, needTx, s)
	printSql(s)
	if deferFn != nil {
		defer deferFn()
	}
	for arg := fn(err); len(arg) != 0; arg = fn(err) {
		if res, err = stmt.Exec(arg...); err == nil {
			var r int64
			if r, err = ResolveResult(res, needId); err == nil {
				ret = append(ret, r)
			}
		}
	}
	return
}

// TxOrNot return an statement, if need transaction, the deferFn will commit or
// rollback transaction, caller should defer or call at last in function
// else only return a normal statement
func TxOrNot(db *sql.DB, needTx bool, s string) (stmt *sql.Stmt, err error, deferFn func()) {
	if needTx {
		var tx *sql.Tx
		tx, err = db.Begin()
		if err == nil {
			stmt, err = tx.Prepare(s)
			deferFn = func() {
				if err == nil {
					tx.Commit()
				} else {
					tx.Rollback()
				}
			}
		}
	} else {
		stmt, err = db.Prepare(s)
	}
	return
}

// ResolveResult resolve sql result, if need id, return last insert id
// else return affected row count
func ResolveResult(res sql.Result, needId bool) (int64, error) {
	if needId {
		return res.LastInsertId()
	} else {
		return res.RowsAffected()
	}
}

// Insert insert model's field to database
// func BatchInsert(db *sql.DB, model []Model, fields []Field, needId bool) ([]int64, error) {
// 	firstModel := model[0]
// 	cols, ph := firstModel.ColsSepPH(fields)
// 	sql := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", firstModel.Table(), cols, ph)
// 	index, indexEnd := 0, len(model)
// 	fn := func(err error) (args []interface{}) {
// 		if err == nil && index < indexEnd {
// 			args = model[index].FieldVals(fields)
// 		}
// 		index++
// 		return
// 	}
// 	return BatchExec(db, sql, fn, true)
// }

// func BatchUpdate(db *sql.DB, model []Model, fields []Field, whereFields []Field) ([]int64, error) {
// 	firstModel := model[0]
// 	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", firstModel.Table(),
// 		firstModel.ColsPH(fields), firstModel.ColsPH(whereFields))
// 	index, indexEnd := 0, len(model)
// 	fn := func(err error) (args []interface{}) {
// 		if err == nil && index < indexEnd {
// 			args = model[index].FieldVals(append(fields, whereFields...))
// 		}
// 		index++
// 		return
// 	}
// 	return BatchExec(db, sql, fn, false)
// }
