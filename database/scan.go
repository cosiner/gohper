package database

import "database/sql"

// Example:
// User should be a Model
// type Users struct {
// 	users  []User
// 	fields uint
// }

// func (u *Users) Make(count int) {
// 	u.users = make([]User, count)
// }

// func (u *Users) Ptrs(index int, ptrs []interface{}) {
// 	u.users[index].Ptrs(u.fields, ptrs)
// }

type Scanner interface {
	// Make will be called twice, first to allocate data space, second to specified
	// the row count
	Make(size int)
	Ptrs(index int, ptrs []interface{})
}

// Scan scan rows to scanner
// if rows has elements, scanner's Make method will be called to allocate space,
// the size will be rowCount, and fields count will get from rows.Columns().
// if therre is no rows, sql.ErrNoRows was returned.
//
// it's mostly designed for that []Model is not you wanted, the []User is.
func Scan(rows *sql.Rows, s Scanner, rowCount int) error {
	index := -1
	var ptrs []interface{}
	defer rows.Close()
	for rows.Next() {
		if index < 0 {
			index = 0
			if cols, err := rows.Columns(); err == nil {
				s.Make(rowCount)
				ptrs = make([]interface{}, len(cols))
			} else {
				return sql.ErrNoRows
			}
		}
		s.Ptrs(index, ptrs)
		if err := rows.Scan(ptrs...); err != nil {
			return err
		}
		index++
	}
	if index < 0 {
		return sql.ErrNoRows
	}
	s.Make(index + 1)
	return nil
}
