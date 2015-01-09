package mgolib

import (
	"fmt"
	"reflect"
	"strings"

	. "github.com/cosiner/golib/errors"
)

type tableMap struct {
	typ     string            // struct type name with pkg, such model.User
	table   string            // corresponding table name of struct
	columns map[string]string // column name mapping
}

func NewTableMap(typeName, tableName string) *tableMap {
	return &tableMap{typeName, tableName, make(map[string]string)}
}

func (tm *tableMap) Columns() map[string]string {
	return tm.columns
}

func (t *tableMap) SetTable(typ reflect.Type, tb string) {
	t.typ = typ.String()
	t.table = tb
}

func (t *tableMap) AddColumn(fieldName string, col string, overwrite bool) error {
	_, has := t.columns[fieldName]
	if has && !overwrite {
		return Err("Column exist")
	}

	t.columns[fieldName] = col
	return nil
}

// map a  struct to database table
// table name is struct name in the form of snake if not define in first field's tag, such as `table:"user"`
// else use the name defined in field tag, struct field is same as table, if no tag defined, use field name
func MapStructToTable(t reflect.Type, scopeSep string) (*tableMap, error) {
	if !StructKind(t) {
		return nil, fmt.Errorf("%s is not struct", t.Name())
	}
	tablename, hasfield := parseTableNameFromStruct(t)
	if !hasfield {
		return nil, fmt.Errorf("%s has no field", t.Name())
	}

	tm := newTableMap(t.String(), tablename)

	pq := &ParseQueue{}
	pq.Enqueue(nil, "", false, t)
	for len(*pq) > 0 {
		st, _ := pq.Dequeue()
		parseStruct(st, pq, tm)
	}

	return tm, nil
}

func parseStruct(st *StructType, pq *ParseQueue, tm *tableMap) {
	if !StructKind(st) {
		return
	}
	for i, n := 0, st.NumField(); i < n; i++ {
		f := st.Field(i)
		tag := f.Tag

		if disabled(tag) || UnExported(f) {
			continue
		}

		if StructKind(f.Type) {
			pq.Enqueue(st, f.Name, f.Anonymous, f.Type)
		} else {
			if st.allprefix == nil {
				st.PrefixCombinations(FIELD_SCOPE_SEPRATOR)
			}
			columnName := parseColumnNameFromField(f)
			fmt.Println(st.allprefix)
			for _, prefix := range st.allprefix {
				tm.addColumn(prefix, columnName, false)
			}
		}
	}
}

func disabled(tag reflect.StructTag) bool {
	return string(tag) == COLUMN_DISABLED || tag.Get(COLUMN_PREFIX) == COLUMN_DISABLED
}

func tableName(t reflect.Type) (name string, hasfield bool) {
	t = RealType(t)
	if t.NumField() > 0 {
		name = t.Field(0).Tag.Get(TABLE_PREFIX)
		hasfield = true
	}
	if name == "" {
		name = strings.ToLower(t.Name())
	}
	return SnakeString(name), hasfield
}

func columnName(f reflect.StructField) (name string) {
	if name = f.Tag.Get(COLUMN_PREFIX); name == "" {
		name = f.Name
	}
	return SnakeString(name)
}

// Represent a struct
type StructType struct {
	names        []string // field name, include parents, for embed struct field non-anonymous
	anonymouses  []bool   // whether or not struct field is anonymous, if true, must use prefix
	allprefix    []string
	reflect.Type // struct to parse
}

// all the accessed prefix for struct field calculate by parents' name and anonymous attribute
func (st *StructType) PrefixCombinations(sep string) {
	//TODO: every struct can reuse parent's allprefix, only to combine own's name
	all := make([]string, 0)
	for i, name := range st.names {
		all = AppendEach(all, sep, name, st.anonymouses[i]) // if anonymous, copy
	}
	fmt.Println(len(all))
	st.allprefix = all
	for _, str := range all {
		fmt.Println(str)
	}
}

// apend sep+str to each of all, if copy, return old and new, else, append after origin
func AppendEach(all []string, sep, str string, optional bool) []string {
	target := all
	if len(all) == 0 {
		target = append(target, str)
		if optional {
			target = append(target, "")
		}
	} else {
		if optional {
			target = make([]string, len(all))
		}
		for i, n := 0, len(all); i < n; i++ {
			target[i] = str
			if all[i] != "" {
				target[i] = all[i] + sep
			}
		}
		if optional {
			target = append(target, all...)
		}
	}
	return target
}

type ParseQueue []*StructType //Parse Queue

func (pq *ParseQueue) Enqueue(parent *StructType, prefix string, anonymous bool, t reflect.Type) *StructType {
	names := make([]string, 0)
	anonymouses := make([]bool, 0)
	if parent != nil {
		names = parent.names
		anonymouses = parent.anonymouses
	}
	names = append(names, prefix)
	anonymouses = append(anonymouses, anonymous)
	st := &StructType{names, anonymouses, nil, t}
	*pq = append(*pq, st)
	return st
}

func (pq *ParseQueue) Dequeue() (st *StructType, has bool) {
	if len(*pq) > 0 {
		st = (*pq)[0]
		*pq = (*pq)[1:]
		has = true
	}
	return
}
