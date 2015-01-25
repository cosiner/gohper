package redis

import (
	"time"

	. "github.com/cosiner/golib/errors"

	"github.com/cosiner/gomodule/config"

	"github.com/garyburd/redigo/redis"
)

var (
	ToString = redis.String
	ToBytes  = redis.Bytes
	ToInt    = redis.Int
)

// Redis store
type RedisStore struct {
	connPool *redis.Pool // redis connection pool
}

func NewRedisStore(conf string) (*RedisStore, error) {
	rc := new(RedisStore)
	err := rc.Init(conf)
	return rc, err
}

func (rc *RedisStore) Init(conf string) error {
	c := config.NewConfig(config.LINE)
	c.ParseString(conf)
	maxidle := c.IntValDef("maxidle", 3)
	idleTimeout := c.IntValDef("idletimeout", 180)
	if maxidle < 0 || idleTimeout < 0 {
		return Err("Wrong format for redis cache")
	}
	addr := c.ValDef("addr", "")
	if addr == "" {
		return Err("Not specified redis server address")
	}
	rc.connPool = &redis.Pool{
		MaxIdle:     maxidle,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (c redis.Conn, err error) {
			c, err = redis.Dial("tcp", addr)
			return
		},
	}
	return nil
}

func (rc *RedisStore) Query(cmd string, args ...interface{}) (reply interface{}, err error) {
	c := rc.connPool.Get()
	reply, err = c.Do(cmd, args...)
	c.Close()
	return
}

func (rc *RedisStore) Update(cmd string, args ...interface{}) error {
	c := rc.connPool.Get()
	_, err := c.Do(cmd, args...)
	c.Close()
	return err
}

func (rc *RedisStore) IsExist(key string) (bool, error) {
	return redis.Bool(rc.Query("EXISTS", key))
}

func (rc *RedisStore) IsHExist(h, key string) (bool, error) {
	return redis.Bool(rc.Query("HEXISTS", h, key))
}

func (rc *RedisStore) Get(key string) (interface{}, error) {
	return rc.Query("GET", key)
}

func (rc *RedisStore) HGet(h, key string) (interface{}, error) {
	return rc.Query("HGET", h, key)
}

func (rc *RedisStore) Set(key string, val interface{}) error {
	return rc.Update("SET", key, val)
}

func (rc *RedisStore) HSet(h, key string, val interface{}) error {
	return rc.Update("HSET", h, key, val)
}

func (rc *RedisStore) SetWithExpire(key string, val interface{}, timeout uint64) error {
	return rc.Update("SETEX", key, timeout, val)
}

func (rc *RedisStore) SetExpire(key string, timeout uint64) error {
	return rc.Update("SETEX", key, timeout)
}

func (rc *RedisStore) Modify(key string, val interface{}) (success bool, err error) {
	if success, err = rc.IsExist(key); err == nil && success {
		err = rc.Set(key, val)
	}
	return
}

func (rc *RedisStore) Remove(key string) error {
	return rc.Update("DEL", key)
}

func (rc *RedisStore) Incr(key string) error {
	return rc.Update("INCR", key)
}

func (rc *RedisStore) Decr(key string) error {
	return rc.Update("DECR", key)
}
