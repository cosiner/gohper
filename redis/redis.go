package redis

import (
	"time"

	. "github.com/cosiner/golib/errors"

	"github.com/cosiner/gomodule/config"

	"github.com/garyburd/redigo/redis"
)

// To* is a set of function to convert redis replay to wanted type
var (
	ToString  = redis.String
	ToBytes   = redis.Bytes
	ToInt     = redis.Int
	ToInt64   = redis.Int64
	ToUint64  = redis.Uint64
	ToFloat64 = redis.Float64
	ToBool    = redis.Bool
)

// RedisStore is a wrapper for "github.com/garyburd/redigo/redis"
type RedisStore struct {
	connPool *redis.Pool // redis connection pool
}

// NewRedisStore return an new RedisStore with given conf
// conf like maxidle=*&idletimeout=*&addr=*
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

// Query excute operation that need an replay value
func (rc *RedisStore) Query(cmd string, args ...interface{}) (reply interface{}, err error) {
	c := rc.connPool.Get()
	reply, err = c.Do(cmd, args...)
	c.Close()
	return
}

// Update execute an operation that without replay
func (rc *RedisStore) Update(cmd string, args ...interface{}) error {
	c := rc.connPool.Get()
	_, err := c.Do(cmd, args...)
	c.Close()
	return err
}

// IsExist check whether given key exist in redis server
func (rc *RedisStore) IsExist(key string) (bool, error) {
	return redis.Bool(rc.Query("EXISTS", key))
}

// IsHExist check whether given key exist in given hash table in redis server
func (rc *RedisStore) IsHExist(h, key string) (bool, error) {
	return redis.Bool(rc.Query("HEXISTS", h, key))
}

// Get return value bind to given key
func (rc *RedisStore) Get(key string) (interface{}, error) {
	return rc.Query("GET", key)
}

// HGet return value bind to given key in given hash table
func (rc *RedisStore) HGet(h, key string) (interface{}, error) {
	return rc.Query("HGET", h, key)
}

// Set bind an value to key
func (rc *RedisStore) Set(key string, val interface{}) error {
	return rc.Update("SET", key, val)
}

// HSet bind an value to key in given hash table
func (rc *RedisStore) HSet(h, key string, val interface{}) error {
	return rc.Update("HSET", h, key, val)
}

// SetWithExpire bind an value to key and set it's expire time
func (rc *RedisStore) SetWithExpire(key string, val interface{}, lifetime int64) error {
	return rc.Update("SETEX", key, lifetime, val)
}

// SetExpire set expire time for exist key
func (rc *RedisStore) SetExpire(key string, lifetime int64) error {
	return rc.Update("SETEX", key, lifetime)
}

// Modify only update exist key's binding value to new value, if update success return true
func (rc *RedisStore) Modify(key string, val interface{}) (success bool, err error) {
	if success, err = rc.IsExist(key); err == nil && success {
		err = rc.Set(key, val)
	}
	return
}

// Remove remove exist key
func (rc *RedisStore) Remove(key string) error {
	return rc.Update("DEL", key)
}

// HRemove remove exist key from an hash table
func (rc *RedisStore) HRemove(h, key string) error {
	return rc.Update("HDEL", h, key)
}

// Incr increase value bind to given key
func (rc *RedisStore) Incr(key string) error {
	return rc.Update("INCR", key)
}

// Decr decrease value bind to given key
func (rc *RedisStore) Decr(key string) error {
	return rc.Update("DECR", key)
}

// Destroy destroy redis store
func (rc *RedisStore) Destroy() {
	rc.connPool.Close()
}
