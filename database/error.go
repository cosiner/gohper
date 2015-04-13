package database

import (
	"database/sql"
	"strings"

	"github.com/cosiner/gohper/lib/types"
)

func ErrForDuplicateKey(err error, newErrFunc func(key string) error) error {
	const duplicate = "Duplicate"
	const forKey = "for key"
	if err != nil {
		s := err.Error()
		index := strings.Index(s, duplicate)
		if index >= 0 {
			s = s[index+len(duplicate):]
			if index = strings.Index(s, forKey) + len(forKey); index >= 0 {
				s, _ = types.TrimQuote(s[index:])
				if e := newErrFunc(s); e != nil {
					err = e
				}
			}
		}
	}
	return err
}

func ErrForNoRows(err, newErr error) error {
	if err == sql.ErrNoRows {
		err = newErr
	}
	return err
}
