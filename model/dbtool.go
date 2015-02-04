package model

import (
	"database/sql"
	"fmt"
)

type (
	fieldsSQLCache map[uint]string
	modelSQLCache  map[string]fieldsSQLCache
	newSQLFunc     func(Model, []Field, []Field) string
)

var (
	insertSQLCache modelSQLCache = make(modelSQLCache)
	updateSQLCache modelSQLCache = make(modelSQLCache)
	deleteSQLCache modelSQLCache = make(modelSQLCache)
	selectSQLCache modelSQLCache = make(modelSQLCache)

	printSQL = func(bool, string) {}
)

// EnableSqlPrint enable print sql
func EnableSqlPrint() {
	printSQL = func(fromCache bool, sql string) {
		fmt.Printf("fromcache:%s, sql:%s\n", fromCache, sql)
	}
}

// fieldsSig return fields signature as fieldset
func fieldsSig(fields []Field) uint {
	return NewFieldSet(fields...).Uint()
}

// cacheGet get sql from cache, if not exist in cache, use newSQL to create a new one
func cacheGet(modelCache modelSQLCache, model Model, fields, whereFields []Field,
	newSQL newSQLFunc) (sql string) {
	table := model.Table()
	fieldsCache := modelCache[table]
	sig := fieldsSig(fields) << model.FieldCount()
	if len(whereFields) != 0 {
		sig |= fieldsSig(whereFields)
	}
	if sig == 0 {
		return ""
	}
	var has bool
	if has = (fieldsCache != nil); has {
		sql, has = fieldsCache[sig]
	} else {
		fieldsCache = make(fieldsSQLCache)
		modelCache[table] = fieldsCache
	}
	if !has {
		sql = newSQL(model, fields, whereFields)
		fieldsCache[sig] = sql
	}
	return sql
}

// insert
func sqlForInsert(model Model, fields, _ []Field) string {
	cols, ph := model.ColsSepPH(fields)
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", model.Table(), cols, ph)
}

// Insert insert model's field to database
func Insert(db *sql.DB, model Model, fields []Field, needId bool) (int64, error) {
	sql := cacheGet(insertSQLCache, model, fields, nil, sqlForInsert)
	args := model.FieldVals(fields)
	return Exec(db, sql, args, needId)
}

func sqlforUpdate(model Model, fields, whereFields []Field) string {
	return fmt.Sprintf("UPDATE %s SET %s WHERE %s", model.Table(),
		model.ColsPH(fields), model.ColsPH(whereFields))
}

// Update update model
func Update(db *sql.DB, model Model, fields, whereFields []Field) (int64, error) {
	sql := cacheGet(updateSQLCache, model, fields, whereFields, sqlforUpdate)
	args := model.FieldVals(append(fields, whereFields...))
	return Exec(db, sql, args, false)
}

func sqlForDelete(model Model, _, whereFields []Field) string {
	return fmt.Sprintf("DELETE FROM %s WHERE %s", model.Table(), model.ColsPH(whereFields))
}

// Delete delete model
func Delete(db *sql.DB, model Model, whereFields []Field) (int64, error) {
	sql := cacheGet(deleteSQLCache, model, nil, whereFields, sqlForDelete)
	args := model.FieldVals(whereFields)
	return Exec(db, sql, args, false)
}

func sqlForSelect(model Model, fields, whereFields []Field) string {
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s", model.Cols(fields), model.Table(),
		model.ColsPH(whereFields))
}

// SelectOne select model from database
func SelectOne(db *sql.DB, model Model, fields, whereFields []Field) error {
	sql := cacheGet(selectSQLCache, model, fields, whereFields, sqlForSelect)
	args := model.FieldVals(whereFields)
	row := db.QueryRow(sql, args...)
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
