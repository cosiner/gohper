package mgolib

/*
import (
	"reflect"
	"sync"
)

type TableMapper struct {
	tables map[string]*tableMap
	lock   sync.Mutex
}

func NewTableMapper() *TableMapper {
	m := &TableMapper{}
	m.tables = make(map[string]*tableMap)
	return m
}

func (t *TableMapper) TableMapping(obj interface{}) *tableMap {
	return t.tables[reflect.TypeOf(obj).String()]
}

func (m *TableMapper) ParseFromValue(objs ...interface{}) (err error) {
	for _, o := range objs {
		t := RealType(reflect.TypeOf(o))
		_, err = m.ParseFromType(t, true)
		if err != nil {
			return
		}
	}
	return
}

type TableMapFunc func() (*tableMap, error)

func (m *TableMapper) ParseFromFunc(fn TableMapFunc, overwrite bool) (tm *tableMap) {
	if tm, err := fn(); err == nil {
		m.AddTableMap(tm, overwrite)
	}
	return
}

func (m *TableMapper) ParseFromType(t reflect.Type, overwrite bool) (tm *tableMap, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	var has bool
	t = RealType(t)
	if tm, has = m.tables[t.String()]; has && !overwrite {
		return
	}
	if tm, err = MapStructToTable(t, FIELD_SCOPE_SEPRATOR); err == nil {
		m.tables[t.String()] = tm
	}
	return
}

func (m *TableMapper) AddTableMap(tm *tableMap, overwrite bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if overwrite {
		m.tables[tm.typ] = tm
	} else {
		if _, has := m.tables[tm.typ]; !has {
			m.tables[tm.typ] = tm
		}
	}
}

func (m *TableMapper) TableName(t reflect.Type) (table string, err error) {
	tm, err := m.ParseFromType(t, false)
	if err == nil {
		table = tm.table
	}
	return
}

func (m *TableMapper) ColumnName(t reflect.Type, fieldName string) (column string) {
	tm, err := m.ParseFromType(t, false)
	if err == nil {
		col, _ := tm.columns[fieldName]
		column = col
	}
	return
}

var (
	TABLE_PREFIX         = "table"
	COLUMN_PREFIX        = "column"
	COLUMN_DISABLED      = "-"
	FIELD_SCOPE_SEPRATOR = "."
)
*/
