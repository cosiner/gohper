package conv

import (
	"strconv"
	"time"

	"github.com/cosiner/gohper/errors"
)

const (
	ErrKeyNotFound = errors.Err("key not found")
	ErrBadFormat   = errors.Err("")
)

type Values struct {
	Timefunc func(string, string) (time.Time, error)
	Timefmt  string

	Err  error
	Vals map[string]string
}

func (v *Values) String(name string) (string, error) {
	val, has := v.Vals[name]
	if !has {
		return "", errors.Newf("%s is not found", name)
	}
	return val, nil
}

func (v *Values) Int(name string) (int, error) {
	val, err := v.String(name)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(val)
}

func (v *Values) Int64(name string) (int64, error) {
	val, err := v.String(name)
	if err != nil {
		return 0, err
	}

	return Atoi64(val)
}

func (v *Values) Uint(name string) (uint, error) {
	val, err := v.String(name)
	if err != nil {
		return 0, err
	}

	return Atou(val)
}

func (v *Values) Uint64(name string) (uint64, error) {
	val, err := v.String(name)
	if err != nil {
		return 0, err
	}
	return Atou64(val)
}

func (v *Values) Bool(name string) (bool, error) {
	val, err := v.String(name)
	if err != nil {
		return false, err
	}
	return Atob(val)
}

func (v *Values) Float64(name string) (float64, error) {
	val, err := v.String(name)
	if err != nil {
		return 0, err
	}
	return Atof(val)
}

func Atoi64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func Atou(s string) (uint, error) {
	u64, err := strconv.ParseUint(s, 10, 0)
	return uint(u64), err
}

func Atou64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

func Atob(s string) (bool, error) {
	return strconv.ParseBool(s)
}

func Atof(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func I64toa(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Utoa(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}

func U64toa(u uint64) string {
	return strconv.FormatUint(u, 10)
}

func Btoa(b bool) string {
	return strconv.FormatBool(b)
}

func Ftoa(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}
