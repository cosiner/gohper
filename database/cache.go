package database

import (
	"database/sql"
	"fmt"
)

type (
	// cacheItem keeps the sql and prepared statement of it
	cacheItem struct {
		sql  string
		stmt *sql.Stmt
	}

	// Cacher store all the sql, statement by sql type and id
	// typically, the sql id of predefied sql type is
	// fields << numField( of Model) | whereFields,
	// it's used for a single model
	//
	// if custom is necessary, call cache.ExtendSQLType(cache.Types()+1) to make
	// a new type, the id is bring your owns, also you can still use the standard
	// FieldIdentity(fields, whereFields) if possible
	Cacher struct {
		cache []map[uint]cacheItem // [typeCount]map[id]{sql, stmt}
		db    *DB
	}
)

const (
	// These are five predefined sql types
	INSERT uint = iota
	DELETE
	UPDATE
	LIMIT_SELECT
	SELECT_ONE
	defaultTypeEnd
)

var (
	// SQLTypes defines the default sql types count, it's default applied to all
	// models.
	// Change it before register any models.
	SQLTypes = defaultTypeEnd
	printSQL = func(_ bool, _ string) {}
)

// SQLPrint enable sql print for each operation
func SQLPrint(enable bool, formatter func(formart string, v ...interface{})) {
	if enable {
		if formatter == nil {
			formatter = func(format string, v ...interface{}) {
				fmt.Printf(format, v...)
			}
		}
		printSQL = func(fromcache bool, sql string) {
			formatter("[SQL]CachedSQL:%t, sql:%s\n", fromcache, sql)
		}
	}
}

// NewCacher create a common sql and statement cacher with given types count
// this will make no parameter checks
//
// the DB instance and each TypeInfo already embed a Cacher, typically, it's not
// necessary to call this
func NewCacher(types uint, db *DB) *Cacher {
	c := &Cacher{
		cache: make([]map[uint]cacheItem, types),
		db:    db,
	}
	for i := uint(0); i < types; i++ {
		c.cache[i] = make(map[uint]cacheItem)
	}
	return c
}

// ExtendSQLType typically used to extend types of Cacher, but it also can be used
// to shrink the cacher, return value will be the new types count you passed in
//
// Example:
// //a.go
// newType1 := c.ExtendSQLType(c.SQLTypes()+1)
// //b.go
// newType2 := c.ExtendSQLType(c.SQLTypes()+1)
func (c *Cacher) ExtendSQLType(typ uint) uint {
	if l := uint(len(c.cache)); typ > l {
		cache := make([]map[uint]cacheItem, typ)
		copy(cache, c.cache)
		for ; l < typ; l++ {
			cache[l] = make(map[uint]cacheItem)
		}
		c.cache = cache
	} else {
		c.cache = c.cache[:typ]
	}
	return typ - 1
}

// StmtById search a prepared statement for given sql type by id, if not found,
// create with the creator, and prepared the sql to a statement, cache it, then
// return
func (c *Cacher) StmtById(typ, id uint, create func() string) (*sql.Stmt, error) {
	if item, has := c.cache[typ][id]; has {
		printSQL(true, item.sql)
		return item.stmt, nil
	}
	sql_ := create()
	stmt, err := c.db.Prepare(sql_)
	if err == nil {
		c.cache[typ][id] = cacheItem{sql: sql_, stmt: stmt}
		printSQL(false, sql_)
	}
	return stmt, err
}

// SQLTypes return the sql types count of current Cacher
func (c *Cacher) SQLTypes() uint {
	return uint(len(c.cache))
}

// GetStmt get sql and statement from cacher, if not found, "" and nil was returned
func (c *Cacher) GetStmt(typ, id uint) (string, *sql.Stmt) {
	if item, has := c.cache[typ][id]; has {
		return item.sql, item.stmt
	}
	return "", nil
}

// SetStmt prepare a sql to statement, cache then return it
func (c *Cacher) SetStmt(typ uint, id uint, sql string) (*sql.Stmt, error) {
	stmt, err := c.db.Prepare(sql)
	if err == nil {
		c.cache[typ][id] = cacheItem{
			sql:  sql,
			stmt: stmt,
		}
	}
	return stmt, err
}
