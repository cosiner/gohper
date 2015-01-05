package mgolib

import (
	"bytes"
	"strings"
)

var Tables = make(map[string]string)

type (
	IdType  int32
	EncType [64]byte
)

type Model interface {
	Table() string
	Column(name string) string
}

type SQLRunner interface {
}

type InsertRunner struct {
	tab  string
	cols []string
}

type QueryRunner struct {
	tabs   []string // tables
	cols   []string // columns
	where  []*Cond  // condition
	having []*Cond
	orders []string
	groups []string
	sql    string
}

func NewQueryRunner() *QueryRunner {
	return &QueryRunner{
		make([]string, 0),
		make([]string, 0),
		make([]*Cond, 0),
		make([]*Cond, 0),
		make([]string, 0),
		make([]string, 0),
	}
}

type Cond struct {
	relation string //
	cond     string
	op       string        // operator
	col      string        // column
	args     []interface{} // arg val
}

const (
	SELECT string = "SELECT "
	FROM   string = " FROM "
	UPDATE string = "UPDATE "
	SET    string = " SET "
	DELETE string = "DELETE "
	INSERT string = "INSERT INTO "
	VALUES string = " VALUES "

	OR  string = " OR "
	AND string = " AND "

	WHERE   string = " WHERE "
	HAVING  string = " HAVING "   // having col1 > n1
	GROUPBY string = " GROUP BY " // group by col1, col2
	ORDERBY string = " ORDER BY " // order by col1, col2
	LIMIT   string = " LIMIT "    // limit n1, n2

	EQ       string = " = "
	NE       string = " <> "
	LT       string = " < "
	GT       string = " > "
	LE       string = " <= "
	GE       string = " >= "
	LIKE     string = " LIKE "    // like %
	IN       string = " IN "      // col1 in (n1, n2)
	BETWEEN  string = " BETWEEN " // col1 between n1 and n2
	NULL     string = " IS NULL "
	NOTNULL  string = " IS NOT NULL "
	ASC      string = " ASC "
	DESC     string = " DESC "
	DISTINCT string = " DISTINCT "

	PLACEHOLDERR string = " ? "
	COMMA        string = ", "

	SQLLENGTH int = 128
)

// join string with seperator by count, such as : str sep str sep....
func JoinString(str, sep string, count int) string {
	if count <= 0 {
		return ""
	}
	s := strings.Repeat(str+sep, count)
	return s[:len(s)-len(sep)]
}

func (c *Cond) String() string {
	switch c.op {
	case EQ, NE, LT, GT, LE, GE, LIKE:
		return c.col + c.op + PLACEHOLDERR
	case LIMIT:
		return LIMIT + JoinString(PLACEHOLDERR, COMMA, len(c.args))
	case IN:
		return c.col + IN + "(" + JoinString(PLACEHOLDERR, COMMA, len(c.args)) + ")"
	case BETWEEN:
		return c.col + BETWEEN + JoinString(PLACEHOLDERR, AND, len(c.args))
	case NULL, NOTNULL:
		return c.col + c.op
	default:
		return ""
	}
}

func colargs(conds []*Cond, argContainer []interface{}) ([]byte, []interface{}) {
	if argContainer == nil {
		argContainer = make([]interface{}, 0)
	}
	buf := bytes.NewBuffer(make([]byte, 0))
	for _, cond := range conds {
		if buf.Len() != 0 && cond.op != LIMIT {
			buf.WriteString(cond.relation)
		}
		buf.WriteString(cond.String())
		argContainer = append(argContainer, cond.args...)
	}
	return buf.Bytes(), argContainer
}

func (qr *QueryRunner) Sql() (string, []interface{}) {
	sqlwriter := bytes.NewBuffer(make([]byte, SQLLENGTH))
	sqlargs := make([]interface{}, 0)
	sqlwriter.WriteString(SELECT)
	sqlwriter.WriteString(strings.Join(qr.cols, COMMA))
	sqlwriter.WriteString(FROM)
	sqlwriter.WriteString(strings.Join(qr.tabs, COMMA))

	where, sqlargs := colargs(qr.where, sqlargs)
	if len(where) > 0 {
		sqlwriter.WriteString(WHERE)
		sqlwriter.Write(where)
	}
	if len(qr.groups) > 0 {
		sqlwriter.WriteString(GROUPBY)
		sqlwriter.WriteString(strings.Join(qr.groups, COMMA))
	}
	having, sqlargs := colargs(qr.having, sqlargs)
	if len(having) > 0 {
		sqlwriter.WriteString(HAVING)
		sqlwriter.Write(having)
	}
	if len(qr.orders) > 0 {
		sqlwriter.WriteString(ORDERBY)
		sqlwriter.WriteString(strings.Join(qr.orders, ","))
	}

	return sqlwriter.String(), sqlargs
}

func (qr *QueryRunner) From(tables ...string) *QueryRunner {
	qr.tabs = append(qr.tabs, tables...)
	return qr
}

func (qr *QueryRunner) Select(cols ...string) *QueryRunner {
	qr.cols = append(qr.cols, cols...)
	return qr
}

func Distinct(cols ...string) string {
	return strings.Join(cols, ",")
}

func Alias(col, alias string) string {
	return col + " " + alias
}

func Asc(col string) string {
	return col + ASC
}
func Desc(col string) string {
	return col + DESC
}

func (qr *QueryRunner) Where(relation string, op string, col string, args ...interface{}) *QueryRunner {
	c := &Cond{relation, WHERE, op, col, args}
	qr.where = append(qr.where, c)
	return qr
}

func (qr *QueryRunner) Having(relation string, op string, col string, args ...interface{}) *QueryRunner {
	c := &Cond{relation, HAVING, op, col, args}
	qr.having = append(qr.having, c)
	return qr
}

func (qr *QueryRunner) Orderby(cols ...string) *QueryRunner {
	qr.orders = append(qr.orders, cols...)
	return qr
}

func (qr *QueryRunner) Groupby(cols ...string) *QueryRunner {
	qr.groups = append(qr.groups, cols...)
	return qr
}
