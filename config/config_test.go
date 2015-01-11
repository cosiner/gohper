package config

import (
	"testing"

	"github.com/cosiner/golib/test"
)

func TestConf(t *testing.T) {
	cfg := NewConfig(INI)
	cfg.ParseString(`
aa=bb
[db]
driver=mysql
host=localhost
port=3306
user=root
password=root
database=test
config=charset=utf-8
pool_max_open=100
pool_max_idle=20
		`)
	val, _ := cfg.Val("aa")
	test.AssertEq(t, "C1", "bb", val)

	_, has := cfg.Val("driver")
	test.AssertFalse(t, "C2", has)

	test.AssertEq(t, "C3", cfg.DefSec(), cfg.CurrSec())

	cfg.SetCurrSec("db")
	test.AssertEq(t, "C4", 3306, cfg.IntVal("port", 0))
}

func BenchmarkConf(b *testing.B) {
	cfg := NewConfig(INI)
	for i := 0; i < b.N; i++ {
		cfg.ParseString(`
aa=bb
[db]
driver=mysql
host=localhost
port=3306
user=root
password=root
database=test
config=charset=utf-8
pool_max_open=100
pool_max_idle=20aa=bb
[dba]
driver=mysql
host=localhost
port=3306
user=root
password=root
database=test
config=charset=utf-8
pool_max_open=100
pool_max_idle=20aa=bb
[dbd]
driver=mysql
host=localhost
port=3306
user=root
password=root
database=test
config=charset=utf-8
pool_max_open=100
pool_max_idle=20aa=bb
[dbz]
driver=mysql
host=localhost
port=3306
user=root
password=root
database=test
config=charset=utf-8
pool_max_open=100
pool_max_idle=20
[dbd]
driver=mysql
host=localhost
port=3306
user=root
password=root
database=test
config=charset=utf-8
pool_max_open=100
pool_max_idle=20aa=bb
[dbz]
driver=mysql
host=localhost
port=3306
user=root
password=root
database=test
config=charset=utf-8
pool_max_open=100
pool_max_idle=20
[dbd]
driver=mysql
host=localhost
port=3306
user=root
password=root
database=test
config=charset=utf-8
pool_max_open=100
pool_max_idle=20aa=bb
[dbz]
driver=mysql
host=localhost
port=3306
user=root
password=root
database=test
config=charset=utf-8
pool_max_open=100
pool_max_idle=20
		`)
	}
}
