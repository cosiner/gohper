package database

import (
	"database/sql"
	"fmt"

	"github.com/cosiner/golib/types"

	. "github.com/cosiner/golib/errors"
)

type (
	// TypeInfo represent information of type
	// contains field count, table name, field names, field offsets
	TypeInfo struct {
		NumField uint
		Table    string
		Fields   []string
	}

	// Model represent a database model
	Model interface {
		Table() string
		FieldValues(*types.LightBitSet) []interface{}
		FieldPtrs(*types.LightBitSet) []interface{}
		New() Model
	}

	// DB holds database connection, all typeinfos, and sql cache
	DB struct {
		*sql.DB
		types          map[string]*TypeInfo
		InsertSQLCache SQLCache
		UpdateSQLCache SQLCache
		DeleteSQLCache SQLCache
		SelectSQLCache SQLCache
	}

	// SQLCache is container of cached sql
	SQLCache map[uint]string

	// sqlCreatorFunc is type of sql creator
	sqlCreatorFunc func(ti *TypeInfo, fields, whereFields *types.LightBitSet) string
)

const (
	// _FIELD_SEP is seperator of columns
	_FIELD_SEP = ","
	// _FIELD_TAG is tag name of database column
	_FIELD_TAG = "column"
)

var printSQL = func(_ bool, _ string) {}

// Open create a database manager
func Open(driver, dsn string, maxIdle, maxOpen int) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	var dbm *DB
	if err == nil {
		db.SetMaxIdleConns(maxIdle)
		db.SetMaxOpenConns(maxOpen)
		dbm = &DB{
			DB:             db,
			types:          make(map[string]*TypeInfo),
			InsertSQLCache: make(SQLCache),
			UpdateSQLCache: make(SQLCache),
			DeleteSQLCache: make(SQLCache),
			SelectSQLCache: make(SQLCache),
		}
	}
	return dbm, err
}

// TypeInfo return type information of given model
// if type info not exist, it will parse it and save type info
func (db *DB) TypeInfo(v Model) *TypeInfo {
	table := v.Table()
	ti := db.types[table]
	if ti == nil {
		ti = db.registerType(v, table)
	}
	return ti
}

// RegisterType register type info, model must not exist
func (db *DB) RegisterType(v Model) {
	table := v.Table()
	Assertf(db.types[table] == nil, "Type %s already registered", table)
	db.registerType(v, table)
}

// registerType save type info of model
func (db *DB) registerType(v Model, table string) *TypeInfo {
	ti := parse(v)
	db.types[table] = ti
	return ti
}

// parse parse model to get type info, model must be struct
// it will first use field tag as column name, if no tag specified,
// use field name's camel_case
func parse(v Model) *TypeInfo {
	typ := types.IndirectType(v)
	fieldNum := typ.NumField()
	fields := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		field := typ.Field(i)
		tag := field.Tag
		fieldName := tag.Get(_FIELD_TAG)
		if fieldName == "" {
			fieldName = field.Name
		}
		fields = append(fields, types.SnakeString(fieldName))
	}
	return &TypeInfo{
		NumField: uint(fieldNum),
		Table:    v.Table(),
		Fields:   fields,
	}
}

// EnableSQLPrint enable sql print for each operation
func EnableSQLPrint() {
	printSQL = func(fromcache bool, sql string) {
		fmt.Printf("fromcache:%t, sql:%s\n", fromcache, sql)
	}
}

// Fields create fieldset from given fields
func Fields(fields ...uint) *types.LightBitSet {
	l := types.NewLightBitSet()
	for _, f := range fields {
		l.Set(f)
	}
	return l
}

// FieldsExcp create fieldset except given fields
func FieldsExcp(numField uint, fields ...uint) *types.LightBitSet {
	set := types.NewLightBitSet()
	set.SetAllBefore(numField)
	for _, f := range fields {
		set.Unset(f)
	}
	return set
}

// CacheGet get sql from cache container, if cache not exist
// then create new sql and save it
func CacheGet(cache SQLCache, sig uint, newSQL func() string) string {
	var sql string
	if sql := cache[sig]; sql == "" {
		sql = newSQL()
		cache[sig] = sql
		printSQL(false, sql)
	} else {
		printSQL(true, sql)
	}
	return sql
}

// cacheGet get sql from cache container, if cache not exist, then create new
// sql and save it
func (db *DB) cacheGet(cache SQLCache, v Model,
	fields, whereFields *types.LightBitSet, newSQL sqlCreatorFunc) string {
	ti := db.TypeInfo(v)
	sig := FieldSignature(ti.NumField, fields, whereFields)
	var sql string
	if sql = cache[sig]; sql == "" {
		sql = newSQL(ti, fields, whereFields)
		cache[sig] = sql
		printSQL(false, sql)
	} else {
		printSQL(true, sql)
	}
	return sql
}

// sqlForInsert create insert sql for given fields
func sqlForInsert(ti *TypeInfo, fields, _ *types.LightBitSet) string {
	count := fields.BitCount()
	cols := ColNames(ti, count, fields)
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", ti.Table,
		ColsString(cols, ""), types.RepeatJoin("?", ",", int(count)))
}

// Insert execure insert operation for model
func (db *DB) Insert(v Model, fields *types.LightBitSet, needId bool) (int64, error) {
	sql := db.cacheGet(db.InsertSQLCache, v, fields, nil, sqlForInsert)
	return db.ExecUpdate(sql, v.FieldValues(fields), needId)
}

// sqlForInsert create update sql for given fields
func (db *DB) sqlForUpdate(ti *TypeInfo, fields, whereFields *types.LightBitSet) string {
	cols := ColNames(ti, fields.BitCount(), fields)
	whereCols := ColNames(ti, whereFields.BitCount(), whereFields)
	return fmt.Sprintf("UPDATE %s SET %s WHERE %s", ti.Table, ColsString(cols, "=?"),
		ColsString(whereCols, "=?"))
}

// Insert execure update operation for model
func (db *DB) Update(v Model, fields *types.LightBitSet, whereFields *types.LightBitSet) (int64, error) {
	values := v.FieldValues(fields)
	values2 := v.FieldValues(whereFields)
	newValues := make([]interface{}, 0, len(values)+len(values2))
	copy(newValues, values)
	copy(newValues[len(values):], values2)
	sql := db.cacheGet(db.UpdateSQLCache, v, fields, whereFields, db.sqlForUpdate)
	return db.ExecUpdate(sql, newValues, false)
}

// sqlForInsert create delete sql for given fields
func (db *DB) sqlForDelete(ti *TypeInfo, _, whereFields *types.LightBitSet) string {
	cols := ColNames(ti, whereFields.BitCount(), whereFields)
	return fmt.Sprintf("DELETE FROM %s WHERE %s", ti.Table, ColsString(cols, "=?"))
}

// Insert execure delete operation for model
func (db *DB) Delete(v Model, whereFields *types.LightBitSet) (int64, error) {
	values := v.FieldValues(whereFields)
	sql := db.cacheGet(db.DeleteSQLCache, v, nil, whereFields, db.sqlForUpdate)
	return db.ExecUpdate(sql, values, false)
}

// sqlForInsert create select sql for given fields
func (db *DB) sqlForSelect(ti *TypeInfo, fields, whereFields *types.LightBitSet) string {
	cols := ColNames(ti, fields.BitCount(), fields)
	whereCols := ColNames(ti, whereFields.BitCount(), whereFields)
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s", ColsString(cols, ""), ti.Table,
		ColsString(whereCols, "=?"))
}

// SelectOne select one row from database
func (db *DB) SelectOne(v Model, fields, whereFields *types.LightBitSet) error {
	ptrs := v.FieldPtrs(fields)
	whereValues := v.FieldValues(whereFields)
	sql := db.cacheGet(db.SelectSQLCache, v, fields, whereFields, db.sqlForSelect)
	return db.ExecQueryRow(sql, whereValues, ptrs)
}

// SelectMulti select multiple results from database
func (db *DB) SelectMulti(v Model, fields, whereFields *types.LightBitSet) ([]Model, error) {
	whereValues := v.FieldValues(whereFields)
	sql := db.cacheGet(db.SelectSQLCache, v, fields, whereFields, db.sqlForSelect)
	vs := make([]Model, 0, 5)
	return vs, db.ExecQuery(sql, whereValues, func() []interface{} {
		vs = append(vs, v)
		ptrs := v.FieldPtrs(fields)
		v = v.New()
		return ptrs
	})
}

// ColNames return column names for given fields
func ColNames(ti *TypeInfo, count uint, fields *types.LightBitSet) []string {
	fieldNames := ti.Fields
	names := make([]string, count)
	index := 0
	for i, name := range fieldNames {
		if fields.IsSet(uint(i)) {
			names[index] = name
			index++
		}
	}
	return names[:index]
}

// ColsString join columns with "," and append a suffix to each column
func ColsString(cols []string, suffix string) string {
	return types.SuffixJoin(cols, suffix, _FIELD_SEP)
}

// FieldSignature create signature from fields
func FieldSignature(numField uint, fields, whereFields *types.LightBitSet) uint {
	return fields.Uint()<<numField | whereFields.Uint()
}

// Exec execute a update operation
func (db *DB) ExecUpdate(s string, args []interface{}, needId bool) (ret int64, err error) {
	res, err := db.Exec(s, args...)
	if err == nil {
		ret, err = ResolveResult(res, needId)
	}
	return
}

// ExecQueryRow execute query operation, scan result to given pointers
func (db *DB) ExecQueryRow(s string, args []interface{}, ptrs []interface{}) error {
	row := db.QueryRow(s, args...)
	return row.Scan(ptrs...)
}

// ExecQuery execute query operation
func (db *DB) ExecQuery(s string, args []interface{}, ptrCreator func() []interface{}) error {
	rows, err := db.Query(s, args...)
	if err == nil {
		for rows.Next() {
			if ptrs := ptrCreator(); ptrs != nil {
				if err = rows.Scan(ptrs...); err != nil {
					break
				}
			} else {
				break
			}
		}
		rows.Close()
	}
	return err
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
