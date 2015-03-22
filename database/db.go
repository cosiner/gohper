package database

import (
	"database/sql"
	"fmt"

	"github.com/cosiner/golib/reflect"

	"github.com/cosiner/golib/types"
)

type (
	// TypeInfo represent information of type
	// contains field count, table name, field names, field offsets
	TypeInfo struct {
		NumField    uint
		Table       string
		Fields      []string
		ErrOnNoRows error
	}

	// Model represent a database model
	Model interface {
		Table() string
		FieldValues(*types.LightBitSet) []interface{}
		FieldPtrs(*types.LightBitSet) []interface{}
		ErrorForNoRows() error
		ErrorForDuplicateKey(key string) error
		New() Model
	}

	// DB holds database connection, all typeinfos, and sql cache
	DB struct {
		driver string
		*sql.DB
		types          map[string]*TypeInfo
		InsertSQLCache SQLCache
		UpdateSQLCache SQLCache
		DeleteSQLCache SQLCache
		SelectSQLCache SQLCache
	}

	// SQLCache is container of cached sql
	SQLCache map[uint]string

	// SQLCreatorFunc is type of sql creator
	SQLCreatorFunc func(ti *TypeInfo, fields, whereFields *types.LightBitSet) string
)

const (
	// _FIELD_SEP is seperator of columns
	_FIELD_SEP = ","
	// _FIELD_TAG is tag name of database column
	_FIELD_TAG = "column"

	_MYSQL_DB = "mysql"
)

var (
	printSQL    = func(_ bool, _ string) {}
	EmptyFields = Fields()
)

// EnableSQLPrint enable sql print for each operation
func EnableSQLPrint(enable bool, formatter func(formart string, v ...interface{})) {
	if enable {
		if formatter == nil {
			formatter = func(format string, v ...interface{}) {
				fmt.Printf(format, v...)
			}
		}
		printSQL = func(fromcache bool, sql string) {
			formatter("[SQL]CachedSQL:%t, sql:%s\n", fromcache, sql)
		}
	} else {
		printSQL = func(bool, string) {}
	}
}

// Open create a database manager and connect to database server
func Open(driver, dsn string, maxIdle, maxOpen int) (*DB, error) {
	db := NewDB()
	err := db.Connect(driver, dsn, maxIdle, maxOpen)
	return db, err
}

// NewDB create a new db
func NewDB() *DB {
	return &DB{
		types:          make(map[string]*TypeInfo),
		InsertSQLCache: make(SQLCache),
		UpdateSQLCache: make(SQLCache),
		DeleteSQLCache: make(SQLCache),
		SelectSQLCache: make(SQLCache),
	}
}

// Connect connect to database server
func (db *DB) Connect(driver, dsn string, maxIdle, maxOpen int) error {
	db_, err := sql.Open(driver, dsn)
	if err == nil {
		db_.SetMaxIdleConns(maxIdle)
		db_.SetMaxOpenConns(maxOpen)
		db.driver = driver
		db.DB = db_
	}
	return err
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

// NumField return field number of a model
func (db *DB) NumField(v Model) uint {
	return db.TypeInfo(v).NumField
}

// ColumnName return return column name of a model field
func (db *DB) ColumnName(v Model, field uint) (name string) {
	fields := db.TypeInfo(v).Fields
	if field < uint(len(fields)) {
		name = fields[field]
	}
	return
}

// RegisterType register type info, model must not exist
func (db *DB) RegisterType(v Model) {
	table := v.Table()
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
	typ := reflect.IndirectType(v)
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
	err := sql.ErrNoRows
	if e := v.ErrorForNoRows(); e != nil {
		err = e
	}
	return &TypeInfo{
		NumField:    uint(fieldNum),
		Table:       v.Table(),
		Fields:      fields,
		ErrOnNoRows: err,
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

// CacheGet get sql from cache container, if cache not exist, then create new
// sql and save it
func (db *DB) CacheGet(cache SQLCache, v Model,
	fields, whereFields *types.LightBitSet, newSQL SQLCreatorFunc) (string, *TypeInfo) {
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
	return sql, ti
}

// SQLForInsert create insert sql for given fields
func SQLForInsert(ti *TypeInfo, fields, _ *types.LightBitSet) string {
	cols := ColNames(ti, fields)
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", ti.Table,
		ColsString(cols, ""), types.RepeatJoin("?", ",", int(fields.BitCount())))
}

// Insert execure insert operation for model
func (db *DB) Insert(v Model, fields *types.LightBitSet, needId bool) (int64, error) {
	sql, _ := db.CacheGet(db.InsertSQLCache, v, fields, EmptyFields, SQLForInsert)
	c, err := db.ExecUpdate(sql, v.FieldValues(fields), needId)
	if db.driver == _MYSQL_DB {
		err = ErrForDuplicateKey(err, v.ErrorForDuplicateKey)
	}
	return c, err
}

// SqlForUpdate create update sql for given fields
func SqlForUpdate(ti *TypeInfo, fields, whereFields *types.LightBitSet) string {
	cols := ColNames(ti, fields)
	whereCols := ColNames(ti, whereFields)
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
	sql, ti := db.CacheGet(db.UpdateSQLCache, v, fields, whereFields, SqlForUpdate)
	c, e := db.ExecUpdate(sql, newValues, false)
	return c, ErrForNoRows(e, ti.ErrOnNoRows)
}

// SQLForDelete create delete sql for given fields
func SQLForDelete(ti *TypeInfo, _, whereFields *types.LightBitSet) string {
	cols := ColNames(ti, whereFields)
	return fmt.Sprintf("DELETE FROM %s WHERE %s", ti.Table, ColsString(cols, "=?"))
}

// Insert execure delete operation for model
func (db *DB) Delete(v Model, whereFields *types.LightBitSet) (int64, error) {
	values := v.FieldValues(whereFields)
	sql, ti := db.CacheGet(db.DeleteSQLCache, v, EmptyFields, whereFields, SQLForDelete)
	c, e := db.ExecUpdate(sql, values, false)
	return c, ErrForNoRows(e, ti.ErrOnNoRows)
}

// SQLForSelect create select sql for given fields
func SQLForSelect(ti *TypeInfo, fields, whereFields *types.LightBitSet) string {
	cols := ColNames(ti, fields)
	whereCols := ColNames(ti, whereFields)
	if len(whereCols) == 0 {
		return fmt.Sprintf("SELECT %s FROM %s", ColsString(cols, ""), ti.Table)
	} else {
		return fmt.Sprintf("SELECT %s FROM %s WHERE %s", ColsString(cols, ""),
			ti.Table, ColsString(whereCols, "=?"))
	}
}

// SelectOne select one row from database
func (db *DB) SelectOne(v Model, fields, whereFields *types.LightBitSet) error {
	ptrs := v.FieldPtrs(fields)
	whereValues := v.FieldValues(whereFields)
	sql, ti := db.CacheGet(db.SelectSQLCache, v, fields, whereFields, SQLForSelect)
	return ErrForNoRows(db.ExecQueryRow(sql, whereValues, ptrs), ti.ErrOnNoRows)
}

// SelectMulti select multiple results from database
func (db *DB) SelectMulti(v Model, fields, whereFields *types.LightBitSet) ([]Model, error) {
	whereValues := v.FieldValues(whereFields)
	sql, ti := db.CacheGet(db.SelectSQLCache, v, fields, whereFields, SQLForSelect)
	vs := make([]Model, 0, 5)
	return vs, ErrForNoRows(db.ExecQuery(sql, whereValues, func() []interface{} {
		vs = append(vs, v)
		ptrs := v.FieldPtrs(fields)
		v = v.New()
		return ptrs
	}), ti.ErrOnNoRows)
}

// SQLForCount create select count sql
func SQLForCount(ti *TypeInfo, _, whereFields *types.LightBitSet) string {
	return fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s",
		ti.Table, ColsString(ColNames(ti, whereFields), "=?"))
}

// Count return count of rows for model
func (db *DB) Count(v Model, whereFields *types.LightBitSet) (count uint, err error) {
	sql, _ := db.CacheGet(db.SelectSQLCache, v, EmptyFields, whereFields, SQLForCount)
	err = db.ExecQueryRow(sql, []interface{}{v.FieldValues(whereFields)},
		[]interface{}{&count})
	return
}

// CountWithArgs return count of rows for model use given arguments
func (db *DB) CountWithArgs(v Model, whereFields *types.LightBitSet,
	args []interface{}) (count uint, err error) {
	sql, _ := db.CacheGet(db.SelectSQLCache, v, EmptyFields, whereFields, SQLForCount)
	err = db.ExecQueryRow(sql, args, []interface{}{&count})
	return
}

// ColNames return column names for given fields
func ColNames(ti *TypeInfo, fields *types.LightBitSet) []string {
	fieldNames := ti.Fields
	names := make([]string, fields.BitCount())
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
	sig := fields.Uint() << numField
	sig |= whereFields.Uint()
	return sig
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
			if err = rows.Scan(ptrCreator()...); err != nil {
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
