package db

import (
	"database/sql"
	"sync"
)

const initSqlCacheSize = 10

type DBMaster struct {
	db        *sql.DB
	cachedSql map[string]string
	cacheLock *sync.RWMutex
}

func NewDBMaster() {
	return &DBMaster{
		cachedSql: make(map[string]string),
		cacheLock: new(sync.RWMutex),
	}
}
