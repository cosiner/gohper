package redis

import (
	"time"

	. "github.com/cosiner/gohper/lib/errors"

	"github.com/cosiner/gohper/config"

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

const (
	ErrWrongIdleSetting = Err("Wrong idle settings for redis")
	ErrRedisAddr        = Err("Not specified redis server address")
)

// RedisStore is a wrapper for "github.com/garyburd/redigo/redis"
type RedisStore struct {
	connPool *redis.Pool // redis connection pool
}

// New return an new RedisStore with given conf
// conf like maxidle=*&idletimeout=*&addr=*
func New2(conf string) (*RedisStore, error) {
	rs := new(RedisStore)
	return rs, rs.Init(conf)
}

func New(maxidle, idletimeout int, addr string) (*RedisStore, error) {
	rs := new(RedisStore)
	return rs, rs.init(maxidle, idletimeout, addr)
}

func (rs *RedisStore) Init(conf string) error {
	c := config.NewConfig(config.LINE)
	c.ParseString(conf)
	return rs.init(c.IntValDef("maxidle", 5),
		c.IntValDef("idletimeout", 180),
		c.ValDef("addr", "127.0.0.1:6379"))
}

func (rs *RedisStore) init(maxidle, idleTimeout int, addr string) error {
	if maxidle < 0 || idleTimeout < 0 {
		return ErrWrongIdleSetting
	}
	if addr == "" {
		return ErrRedisAddr
	}
	rs.connPool = &redis.Pool{
		MaxIdle:     maxidle,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (c redis.Conn, err error) {
			c, err = redis.Dial("tcp", addr)
			return
		},
	}
	return nil
}

func (rs *RedisStore) Conn() redis.Conn {
	return rs.connPool.Get()
}

// Query excute operation that need an replay value
func (rs *RedisStore) Query(cmd string, args ...interface{}) (reply interface{}, err error) {
	c := rs.connPool.Get()
	reply, err = c.Do(cmd, args...)
	c.Close()
	return
}

// Update execute an operation that without replay
func (rs *RedisStore) Update(cmd string, args ...interface{}) error {
	c := rs.connPool.Get()
	_, err := c.Do(cmd, args...)
	c.Close()
	return err
}

// IsExist check whether given key exist in redis server
func (rs *RedisStore) IsExist(key string) (bool, error) {
	return redis.Bool(rs.Query("EXISTS", key))
}

// IsHExist check whether given key exist in given hash table in redis server
func (rs *RedisStore) IsHExist(h, key string) (bool, error) {
	return redis.Bool(rs.Query("HEXISTS", h, key))
}

// Get return value bind to given key
func (rs *RedisStore) Get(key string) (interface{}, error) {
	return rs.Query("GET", key)
}

// HGet return value bind to given key in given hash table
func (rs *RedisStore) HGet(h, key string) (interface{}, error) {
	return rs.Query("HGET", h, key)
}

func (rs *RedisStore) HExists(h, key string) (bool, error) {
	return ToBool(rs.Query("HEXISTs", h, key))
}

// Set bind an value to key
func (rs *RedisStore) Set(key string, val interface{}) error {
	return rs.Update("SET", key, val)
}

// HSet bind an value to key in given hash table
func (rs *RedisStore) HSet(h, key string, val interface{}) error {
	return rs.Update("HSET", h, key, val)
}

// SetWithExpire bind an value to key and set it's expire time
func (rs *RedisStore) SetWithExpire(key string, val interface{}, lifetime int64) error {
	return rs.Update("SETEX", key, lifetime, val)
}

// SetExpire set expire time for exist key
func (rs *RedisStore) SetExpire(key string, lifetime int64) error {
	return rs.Update("SETEX", key, lifetime)
}

// Modify only update exist key's binding value to new value, if update success return true
func (rs *RedisStore) Modify(key string, val interface{}) (success bool, err error) {
	if success, err = rs.IsExist(key); err == nil && success {
		err = rs.Set(key, val)
	}
	return
}

// Remove remove exist key
func (rs *RedisStore) Remove(key string) error {
	return rs.Update("DEL", key)
}

func (rs *RedisStore) Exists(key string) (bool, error) {
	return ToBool(rs.Query("EXISTs", key))
}

// HRemove remove exist key from an hash table
func (rs *RedisStore) HRemove(h, key string) error {
	return rs.Update("HDEL", h, key)
}

// Incr increase value bind to given key
func (rs *RedisStore) Incr(key string) error {
	return rs.Update("INCR", key)
}

// Decr decrease value bind to given key
func (rs *RedisStore) Decr(key string) error {
	return rs.Update("DECR", key)
}

// Destroy destroy redis store
func (rs *RedisStore) Destroy() {
	rs.connPool.Close()
}

func (rs *RedisStore) ToString(v interface{}, err error) (string, error) {
	return redis.String(v, err)
}

func (rs *RedisStore) ToBytes(v interface{}, err error) ([]byte, error) {
	return redis.Bytes(v, err)
}

func (rs *RedisStore) ToInt(v interface{}, err error) (int, error) {
	return redis.Int(v, err)
}

func (rs *RedisStore) ToInt64(v interface{}, err error) (int64, error) {
	return redis.Int64(v, err)
}

func (rs *RedisStore) ToUint64(v interface{}, err error) (uint64, error) {
	return redis.Uint64(v, err)
}

func (rs *RedisStore) ToFloat64(v interface{}, err error) (float64, error) {
	return redis.Float64(v, err)
}

func (rs *RedisStore) ToBool(v interface{}, err error) (bool, error) {
	return redis.Bool(v, err)
}
