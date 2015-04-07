package database

import "database/sql"

type (
	// Model represent a database model
	Model interface {
		Table() string
		// FieldValues return all values of the fields
		// and also should reserve some space to store other parameter values
		// it's recommand, not force
		FieldValues(fields uint, reserveSize uint) []interface{}
		FieldPtrs(uint) []interface{}
		New() Model
	}

	// DB holds database connection, all typeinfos, and sql cache
	DB struct {
		// driver string
		*sql.DB
		types map[string]*TypeInfo
		CommonCacher
	}
)

// Open create a database manager and connect to database server
func Open(driver, dsn string, maxIdle, maxOpen int) (*DB, error) {
	db := New()
	err := db.Connect(driver, dsn, maxIdle, maxOpen)
	return db, err
}

// New create a new db
func New() *DB {
	return &DB{
		types:        make(map[string]*TypeInfo),
		CommonCacher: NewCommonCacher(0),
	}
}

// Connect connect to database server
func (db *DB) Connect(driver, dsn string, maxIdle, maxOpen int) error {
	db_, err := sql.Open(driver, dsn)
	if err == nil {
		db_.SetMaxIdleConns(maxIdle)
		db_.SetMaxOpenConns(maxOpen)
		db.DB = db_
	}
	return err
}

// RegisterType register type info, model must not exist
func (db *DB) RegisterType(v Model) {
	table := v.Table()
	db.registerType(v, table)
}

// registerType save type info of model
func (db *DB) registerType(v Model, table string) *TypeInfo {
	ti := parseTypeInfo(v)
	db.types[table] = ti
	return ti
}

func (db *DB) SQLTypeEnd(typ SQLType) {
	db.CommonCacher = db.CommonCacher.SQLTypeEnd(typ)
}

// TypeInfo return type information of given model
// if type info not exist, it will parseTypeInfo it and save type info
func (db *DB) TypeInfo(v Model) *TypeInfo {
	table := v.Table()
	ti := db.types[table]
	if ti == nil {
		ti = db.registerType(v, table)
	}
	return ti
}

// Insert execure insert operation for model
func (db *DB) Insert(v Model, fields uint, needId bool) (int64, error) {
	ti := db.TypeInfo(v)
	sql := ti.CacheGet(INSERT, fields, 0, ti.InsertSQL)
	return db.ExecUpdate(sql, v.FieldValues(fields, 0), needId)
}

// Insert execure update operation for model
func (db *DB) Update(v Model, fields uint, whereFields uint) (int64, error) {
	values := v.FieldValues(fields, 0)
	values2 := v.FieldValues(whereFields, 0)
	newValues := make([]interface{}, len(values)+len(values2))
	copy(newValues, values)
	copy(newValues[len(values):], values2)
	ti := db.TypeInfo(v)
	sql := ti.CacheGet(UPDATE, fields, whereFields, ti.UpdateSQL)
	return db.ExecUpdate(sql, newValues, false)
}

// Insert execure delete operation for model
func (db *DB) Delete(v Model, whereFields uint) (int64, error) {
	values := v.FieldValues(whereFields, 0)
	ti := db.TypeInfo(v)
	sql := ti.CacheGet(DELETE, 0, whereFields, ti.DeleteSQL)
	return db.ExecUpdate(sql, values, false)
}

func (db *DB) limitSelectRows(v Model, fields, whereFields uint, start, count int) (*sql.Rows, error) {
	ti := db.TypeInfo(v)
	sql := ti.CacheGet(LIMIT_SELECT, fields, whereFields, ti.LimitSelectSQL)
	args := append(append(v.FieldValues(whereFields, 2), start), count)
	return db.Query(sql, args...)
}

// SelectOne select one row from database
func (db *DB) SelectOne(v Model, fields, whereFields uint) error {
	rows, err := db.limitSelectRows(v, fields, whereFields, 0, 1)
	if err == nil {
		if rows.Next() {
			err = rows.Scan(v.FieldPtrs(fields)...)
		} else {
			err = sql.ErrNoRows
		}
	}
	rows.Close()
	return err
}

// Select select multiple results from database
func (db *DB) SelectLimit(v Model, fields, whereFields uint, start, count int) (models []Model, err error) {
	rows, err := db.limitSelectRows(v, fields, whereFields, start, count)
	if err == nil {
		has := false
		models = make([]Model, 0, count)
		for rows.Next() {
			has = true
			model := v.New()
			if err = rows.Scan(model.FieldPtrs(fields)...); err != nil {
				models = nil
				break
			} else {
				models = append(models, model)
			}
		}
		if !has {
			err = sql.ErrNoRows
		}
	}
	rows.Close()
	return models, err
}

// Count return count of rows for model
func (db *DB) Count(v Model, whereFields uint) (count uint, err error) {
	return db.CountWithArgs(v, whereFields, v.FieldValues(whereFields, 0))
}

// CountWithArgs return count of rows for model use given arguments
func (db *DB) CountWithArgs(v Model, whereFields uint,
	args []interface{}) (count uint, err error) {
	ti := db.TypeInfo(v)
	sql := ti.CacheGet(LIMIT_SELECT, 0, whereFields, ti.CountSQL)
	rows, err := db.Query(sql, args...)
	if err == nil {
		rows.Next()
		err = rows.Scan(&count)
	}
	return
}

// Exec execute a update operation
func (db *DB) ExecUpdate(s string, args []interface{}, needId bool) (ret int64, err error) {
	res, err := db.Exec(s, args...)
	if err == nil {
		ret, err = ResolveResult(res, needId)
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
