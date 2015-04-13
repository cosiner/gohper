// Package database is a library help for interact with database by model
//
package database

import "database/sql"

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
		db.Cacher = NewCacher(SQLTypes, db)
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
	return db.InsertWith(v, fields, needId, FieldVals(fields, v))
}

func (db *DB) InsertWith(v Model, fields uint, needId bool, args []interface{}) (int64, error) {
	stmt, err := db.TypeInfo(v).InsertStmt(fields)
	return StmtExec(stmt, err, args, needId)
}

func (db *DB) Update(v Model, fields uint, whereFields uint) (int64, error) {
	c1, c2 := FieldCount(fields), FieldCount(whereFields)
	args := make([]interface{}, c1+c2)
	v.Vals(fields, args)
	v.Vals(whereFields, args[c1:])
	return db.UpdateWith(v, fields, whereFields, args)
}

func (db *DB) UpdateWith(v Model, fields uint, whereFields uint, args []interface{}) (int64, error) {
	stmt, err := db.TypeInfo(v).UpdateStmt(fields, whereFields)
	return StmtExec(stmt, err, args, false)
}

func (db *DB) Delete(v Model, whereFields uint) (int64, error) {
	return db.DeleteWith(v, whereFields, FieldVals(whereFields, v))
}

func (db *DB) DeleteWith(v Model, whereFields uint, args []interface{}) (int64, error) {
	stmt, err := db.TypeInfo(v).DeleteStmt(whereFields)
	return StmtExec(stmt, err, args, false)
}

func (db *DB) rows(v Model, fields, whereFields uint, start, count int) (*sql.Rows, error) {
	c := FieldCount(whereFields)
	args := make([]interface{}, c+2)
	v.Vals(whereFields, args)
	args[c], args[c+1] = start, count
	stmt, err := db.TypeInfo(v).LimitSelectStmt(fields, whereFields)
	if err == nil {
		return stmt.Query(args...)
	}
	return nil, err
}

func (db *DB) row(v Model, fields, whereFields uint) (*sql.Rows, error) {
	stmt, err := db.TypeInfo(v).SelectOneStmt(fields, whereFields)
	if err == nil {
		return stmt.Query(FieldVals(whereFields, v)...)
	}
	return nil, err
}

// SelectOne select one row from database
func (db *DB) SelectOne(v Model, fields, whereFields uint) error {
	rows, err := db.row(v, fields, whereFields)
	return ScanOnce(rows, err, FieldPtrs(fields, v))
}

func (db *DB) SelectLimit(v Model, fields, whereFields uint, start, count int) (
	models []Model, err error) {

	rows, err := db.rows(v, fields, whereFields, start, count)
	if err == nil {
		has := false
		for rows.Next() {
			if !has {
				models = make([]Model, 0, count)
				has = true
			}
			model := v.New()
			if err = rows.Scan(FieldPtrs(fields, model)...); err != nil {
				rows.Close()
				return nil, err
			}
			models = append(models, model)
		}
		if !has {
			err = sql.ErrNoRows
		}
		rows.Close()
	}
	return models, err
}

func (db *DB) ScanLimit(v Model, s Scanner, fields, whereFields uint, start, count int) error {
	rows, err := db.rows(v, fields, whereFields, start, count)
	return ScanLimit(rows, err, s, count)
}

// Count return count of rows for model
func (db *DB) Count(v Model, whereFields uint) (count uint, err error) {
	return db.CountWith(v, whereFields, FieldVals(whereFields, v))
}

// CountWith return count of rows for model use given arguments
func (db *DB) CountWith(v Model, whereFields uint,
	args []interface{}) (count uint, err error) {
	ti := db.TypeInfo(v)
	stmt, err := ti.CountStmt(whereFields)
	if err == nil {
		rows, e := stmt.Query(args...)
		err = ScanOnce(rows, e, &count)
	}
	return
}

// ExecUpdate execute a update operation
func (db *DB) ExecUpdate(s string, args []interface{}, needId bool) (ret int64, err error) {
	res, err := db.Exec(s, args...)
	return ResolveResult(res, err, needId)
}

func StmtExec(stmt *sql.Stmt, err error, args []interface{}, needId bool) (int64, error) {
	if err == nil {
		res, err := stmt.Exec(args...)
		return ResolveResult(res, err, needId)
	}
	return 0, err
}

// ResolveResult resolve sql result, if need id, return last insert id
// else return affected row count
func ResolveResult(res sql.Result, err error, needId bool) (int64, error) {
	if err == nil {
		if needId {
			return res.LastInsertId()
		} else {
			return res.RowsAffected()
		}
	}
	return 0, err
}
