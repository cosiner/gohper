package database

import (
	"database/sql"
	"strings"

	"github.com/cosiner/golib/types"
)

func ErrForDuplicateKey(err error, newErrFunc func(key string) error) error {
	const DUPLICATE = "Duplicate"
	s := err.Error()
	index := strings.Index(s, DUPLICATE)
	if index >= 0 {
		s = s[index+len(DUPLICATE):]
		if index = strings.Index(s, "for key") + 7; index >= 0 {
			s, _ = types.TrimQuote(s[index:])
			if e := newErrFunc(s); e != nil {
				return e
			}
		}
	}
	return err
}

func ErrForNoRows(err, newErr error) error {
	if err == sql.ErrNoRows {
		return newErr
	}
	return err
}
