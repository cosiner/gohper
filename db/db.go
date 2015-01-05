package db

import (
	"container/list"
	"database/sql"
	"errors"
	"sync"
)

const initSqlCacheSize = 10

type DBMaster struct {
	db        *sql.DB
	cachedSql map[string]string
	cacheLock *sync.RMutex
}

func NewDBMaster() {
	return &DBMaster{
		cachedSql: make(map[string]string),
		cacheLock: new(sync.RMutex),
	}
}
