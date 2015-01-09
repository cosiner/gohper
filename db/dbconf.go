package db

import (
	"database/sql"
	"fmt"
	"mlib/config"

	. "github.com/cosiner/golib/errors"
)

func OpenDBFromConf(conf *DBConf) (db *sql.DB, err error) {
	if conf == nil {
		return
	}
	dsn, err := DataSourceName(conf)
	if err != nil {
		return
	}
	db, err = sql.Open(conf.Driver, dsn)
	return
}

func OpenDBFromFile(conffile string) *sql.DB {
	conf, err := LoadDBConf(conffile)
	if err != nil {
		panic(err.Error())
	}
	db, err := OpenDB(conf)
	if err != nil {
		panic(err.Error())
	}
	db.SetMaxIdleConns(conf.PoolMaxIdle)
	db.SetMaxOpenConns(conf.PoolMaxOpen)
	return db
}

type DBConf struct {
	Driver      string
	Host        string
	Port        string
	Database    string
	User        string
	Password    string
	Config      string
	PoolMaxOpen int
	PoolMaxIdle int
}

func LoadDBConf(filename string) (c *DBConf, err error) {
	config := config.NewIniConfig()
	config.SetDefSec("db")

	if err = config.ParseFile(filename); err != nil {
		return
	}
	c = &DBConf{}
	c.Driver, _ = config.Val("driver")
	c.Host, _ = config.Val("host")
	c.Port, _ = config.Val("port")
	c.Database, _ = config.Val("database")
	c.User, _ = config.Val("user")
	c.Password, _ = config.Val("password")
	c.Config, _ = config.Val("config")
	c.PoolMaxOpen = config.IntVal("pool_max_open", 0)
	c.PoolMaxIdle = config.IntVal("pool_max_idle", 0)

	return
}

func DataSourceName(conf *DBConf) (string, error) {
	if conf == nil || conf.Database == "" {
		return "", Err("No Config or database is not set")
	}

	switch conf.Driver {
	case "mysql":
		return mysqlDSN(conf), nil
	}
	return "", Err("No support")
}

func mysqlDSN(conf *DBConf) string {
	if conf.Host == "" {
		conf.Host = "localhost"
	}
	if conf.Port == "" {
		conf.Port = "3306"
	}
	if conf.User == "" {
		conf.User = "root"
	}
	if conf.Password != "" {
		conf.Password = ":" + conf.Password
	}

	return fmt.Sprintf("%s%s@tcp(%s:%d)%s?%s",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Database, conf.Config)
}
