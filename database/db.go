package database

import (
	"github.com/cosiner/gohper/lib/types"

	"database/sql"
)

type (
	// Model represent a database model
	Model interface {
		Table() string
		// Vals store values of fields to given slice
		Vals(fields uint, vals []interface{})
		Ptrs(fields uint, ptrs []interface{})
		New() Model
	}

	// DB holds database connection, all typeinfos, and sql cache
	DB struct {
		// driver string
		*sql.DB
		types map[string]*TypeInfo
		*Cacher
	}
)

var FieldCount = types.BitCountUint

// Open create a database manager and connect to database server
func Open(driver, dsn string, maxIdle, maxOpen int) (*DB, error) {
	db := New()
	err := db.Connect(driver, dsn, maxIdle, maxOpen)
	return db, err
}

// New create a new db
func New() *DB {
	return &DB{
		types: make(map[string]*TypeInfo),
	}
}

// Connect connect to database server
func (db *DB) Connect(driver, dsn string, maxIdle, maxOpen int) error {
	db_, err := sql.Open(driver, dsn)
	if err == nil {
		db_.SetMaxIdleConns(maxIdle)
		db_.SetMaxOpenConns(maxOpen)
		db.DB = db_
		db.Cacher = NewCacher(SQLTypeEnd, db)
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
	ti := parseTypeInfo(v, db)
	db.types[table] = ti
	return ti
}

// TypeInfo return type information of given model
// if type info not exist, it will parseTypeInfo it and save type info
func (db *DB) TypeInfo(v Model) *TypeInfo {
	table := v.Table()
	if ti, has := db.types[table]; has {
		return ti
	}
	return db.registerType(v, table)
}

func FieldVals(fields uint, v Model) []interface{} {
	args := make([]interface{}, FieldCount(fields))
	v.Vals(fields, args)
	return args
}

func FieldPtrs(fields uint, v Model) []interface{} {
	ptrs := make([]interface{}, FieldCount(fields))
	v.Ptrs(fields, ptrs)
	return ptrs
}

func (db *DB) Insert(v Model, fields uint, needId bool) (int64, error) {
	ti := db.TypeInfo(v)
	stmt := ti.CacheGet(INSERT, fields, 0, ti.InsertSQL)
	return StmtExec(stmt, FieldVals(fields, v), needId)
}

func (db *DB) Update(v Model, fields uint, whereFields uint) (int64, error) {
	c1, c2 := FieldCount(fields), FieldCount(whereFields)
	args := make([]interface{}, c1+c2)
	v.Vals(fields, args)
	v.Vals(whereFields, args[c1:])
	ti := db.TypeInfo(v)
	stmt := ti.CacheGet(UPDATE, fields, whereFields, ti.UpdateSQL)
	return StmtExec(stmt, args, false)
}

func (db *DB) Delete(v Model, whereFields uint) (int64, error) {
	ti := db.TypeInfo(v)
	stmt := ti.CacheGet(DELETE, 0, whereFields, ti.DeleteSQL)
	return StmtExec(stmt, FieldVals(whereFields, v), false)
}

func (db *DB) limitSelectRows(v Model, fields, whereFields uint, start, count int) (*sql.Rows, error) {
	ti := db.TypeInfo(v)
	stmt := ti.CacheGet(LIMIT_SELECT, fields, whereFields, ti.LimitSelectSQL)
	c := FieldCount(whereFields)
	args := make([]interface{}, c+1)
	v.Vals(whereFields, args)
	args[c], args[c+1] = start, count
	return stmt.Query(args...)
}

// SelectOne select one row from database
func (db *DB) SelectOne(v Model, fields, whereFields uint) error {
	rows, err := db.limitSelectRows(v, fields, whereFields, 0, 1)
	if err == nil {
		if rows.Next() {
			err = rows.Scan(FieldPtrs(fields, v)...)
		} else {
			err = sql.ErrNoRows
		}
	}
	rows.Close()
	return err
}

func (db *DB) SelectLimit(v Model, fields, whereFields uint, start, count int) (
	models []Model, err error) {

	rows, err := db.limitSelectRows(v, fields, whereFields, start, count)
	if err == nil {
		has := false
		for rows.Next() {
			if !has {
				models = make([]Model, 0, count)
				has = true
			}
			model := v.New()
			if err = rows.Scan(FieldPtrs(fields, model)...); err != nil {
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

func (db *DB) ScanLimit(v Model, s Scanner, fields, whereFields uint, start, count int) error {
	rows, err := db.limitSelectRows(v, fields, whereFields, start, count)
	if err == nil {
		err = Scan(rows, s, count)
	}
	return err
}

// Count return count of rows for model
func (db *DB) Count(v Model, whereFields uint) (count uint, err error) {
	return db.CountWithArgs(v, whereFields, FieldVals(whereFields, v))
}

// CountWithArgs return count of rows for model use given arguments
func (db *DB) CountWithArgs(v Model, whereFields uint,
	args []interface{}) (count uint, err error) {
	ti := db.TypeInfo(v)
	stmt := ti.CacheGet(LIMIT_SELECT, 0, whereFields, ti.CountSQL)
	rows, err := stmt.Query(args...)
	if err == nil {
		rows.Next()
		err = rows.Scan(&count)
	}
	return
}

// ExecUpdate execute a update operation
func (db *DB) ExecUpdate(s string, args []interface{}, needId bool) (ret int64, err error) {
	res, err := db.Exec(s, args...)
	if err == nil {
		ret, err = ResolveResult(res, needId)
	}
	return
}

func StmtExec(stmt *sql.Stmt, args []interface{}, needId bool) (ret int64, err error) {
	res, err := stmt.Exec(args...)
	if err == nil {
		ret, err = ResolveResult(res, needId)
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
