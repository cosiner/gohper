package database

import "database/sql"

// // Scanner Example:
// User should be a Model
// type Users struct {
// 	users  []User
// 	fields uint
// }
//
// func (u *Users) Make(count int) {
// 	u.users = make([]User, count)
// }
//
// func (u *Users) Ptrs(index int, ptrs []interface{}) {
// 	u.users[index].Ptrs(u.fields, ptrs)
// }

type Scanner interface {
	// Make will be called twice, first to allocate data space, second to specified
	// the row count
	Make(size int)
	Ptrs(index int, ptrs []interface{})
}

// ScanLimit scan rows to scanner
// if rows has elements, scanner's Make method will be called to allocate space,
// the size will be rowCount, and fields count will get from rows.Columns().
// if therre is no rows, sql.ErrNoRows was returned.
//
// it's mostly designed for that customed search
func ScanLimit(rows *sql.Rows, err error, s Scanner, rowCount int) error {
	if err == nil {
		index := -1
		var ptrs []interface{}
		for rows.Next() {
			if index < 0 {
				index = 0
				cols, _ := rows.Columns()
				s.Make(rowCount)
				ptrs = make([]interface{}, len(cols))
			}
			s.Ptrs(index, ptrs)
			if err = rows.Scan(ptrs...); err != nil {
				rows.Close()
				return err
			}
			index++
		}
		if index < 0 {
			err = sql.ErrNoRows
		} else {
			s.Make(index + 1)
		}
		rows.Close()
	}
	return err
}

// ScanOnce scan once then close rows, if no data, sql.ErrNoRows was returned
func ScanOnce(rows *sql.Rows, err error, ptrs ...interface{}) error {
	if err == nil {
		if rows.Next() {
			err = rows.Scan(ptrs...)
		} else {
			err = sql.ErrNoRows
		}
		rows.Close()
	}
	return err
}
