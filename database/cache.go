package database

import (
	"database/sql"
	"fmt"
)

type (
	cacheItem struct {
		sql  string
		stmt *sql.Stmt
	}

	Cacher struct {
		cache []map[uint]cacheItem // [typeCount]map[id]{sql, stmt}
		db    *sql.DB
	}
)

var printSQL = func(_ bool, _ string) {}

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

// NewCacher create a common sql cacher with given type end,
// if type end is 0, use global type end
func NewCacher(typEnd SQLType, db *DB) *Cacher {
	return &Cacher{
		cache: make([]map[uint]cacheItem, typEnd),
		db:    db.DB,
	}
}

func (c *Cacher) CacheGetById(typ SQLType, id uint, create func() string) *sql.Stmt {
	item, has := c.cache[typ][id]
	if has {
		printSQL(true, item.sql)
		return item.stmt
	} else {
		sql_ := create()
		stmt, err := c.db.Prepare(sql_)
		if err != nil {
			panic(err)
		}
		item = cacheItem{sql: sql_, stmt: stmt}
		c.cache[typ][id] = item
		printSQL(false, sql_)
		return stmt
	}
}

func (c *Cacher) SQLTypeEnd(typ SQLType) SQLType {
	return SQLType(len(c.cache))
}

func (c *Cacher) SetSQLTypeEnd(typ SQLType) {
	if int(typ) > len(c.cache) {
		cache := make([]map[uint]cacheItem, typ)
		copy(cache, c.cache)
		c.cache = cache
	} else {
		c.cache = c.cache[:typ]
	}
}

func (c *Cacher) CacheGet(typ SQLType, id uint) (string, *sql.Stmt) {
	item, has := c.cache[typ][id]
	if has {
		return item.sql, item.stmt
	}
	return "", nil
}

func (c *Cacher) CacheSet(typ SQLType, id uint, sql string) *sql.Stmt {
	stmt, err := c.db.Prepare(sql)
	if err == nil {
		c.cache[typ][id] = cacheItem{
			sql:  sql,
			stmt: stmt,
		}
		return stmt
	}
	panic(err)
}
