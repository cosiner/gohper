package database

import "github.com/cosiner/gohper/lib/types"

type (
	Cols interface {
		String() string
		Paramed() string
		OnlyParam() string
		Length() int
	}
	cols struct {
		cols        []string
		str         string
		paramed     string
		onlyParamed string
	}
	singleCol string
	nilCols   string
)

var (
	FieldCount      = types.BitCountUint
	zeroCols   Cols = nilCols("")
)

// String return columns string join with ",",
// result like "foo, bar"
func (c *cols) String() string {
	if c.str == "" {
		c.str = types.SuffixJoin(c.cols, "", ",")
	}
	return c.str
}

// Paramed return columns string joind with "=?,", last "," was trimed,
// result like "foo=?, bar=?"
func (c *cols) Paramed() string {
	if c.paramed == "" {
		c.paramed = types.SuffixJoin(c.cols, "=?", ",")
	}
	return c.paramed
}

// OnlyParam return columns placeholdered string, each column was replaced with "?"
// result like "?, ?, ?, ?", count of "?" is colums length
func (c *cols) OnlyParam() string {
	if c.onlyParamed == "" {
		c.onlyParamed = types.RepeatJoin("?", ",", len(c.cols))
	}
	return c.onlyParamed
}

func (c *cols) Length() int {
	return len(c.cols)
}

// String return columns string join with ",",
// result like "foo, bar"
func (c singleCol) String() string {
	return string(c)
}

// Paramed return columns string joind with "=?,", last "," was trimed,
// result like "foo=?, bar=?"
func (c singleCol) Paramed() string {
	return string(c) + "=?"
}

// OnlyParam return columns placeholdered string, each column was replaced with "?"
// result like "?, ?, ?, ?", count of "?" is colums length
func (c singleCol) OnlyParam() string {
	return "?"
}

func (c singleCol) Length() int {
	return 1
}

// String return columns string join with ",",
// result like "foo, bar"
func (nilCols) String() string {
	return ""
}

// Paramed return columns string joind with "=?,", last "," was trimed,
// result like "foo=?, bar=?"
func (nilCols) Paramed() string {
	return ""
}

// OnlyParam return columns placeholdered string, each column was replaced with "?"
// result like "?, ?, ?, ?", count of "?" is colums length
func (nilCols) OnlyParam() string {
	return ""
}

func (nilCols) Length() int {
	return 0
}
