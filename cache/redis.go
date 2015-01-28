package cache

import (
	"github.com/cosiner/gomodule/redis"
)

// RedisCache is only a adapter of redis store
type RedisCache struct {
	redisStore *redis.RedisStore
}

func (rc *RedisCache) Init(config string) error {
	rs, err := redis.NewRedisStore(config)
	if err == nil {
		rc.redisStore = rs
	}
	return err
}

func (rc *RedisCache) InitVals(config string, values map[string]interface{}) error {
	rs, err := redis.NewRedisStore(config)
	if err == nil {
		rc.redisStore = rs
		for k, v := range values {
			if err = rs.Set(k, v); err != nil {
				break
			}
		}
	}
	return err
}

// Get by key
func (rc *RedisCache) Get(key string) interface{} {
	v, err := rc.redisStore.Get(key)
	if err != nil {
		v = nil
	}
	return v
}

func (rc *RedisCache) Set(key string, val interface{}) {
	rc.redisStore.Set(key, val)
}

func (rc *RedisCache) Update(key string, val interface{}) bool {
	success, err := rc.redisStore.Modify(key, val)
	if err != nil {
		success = false
	}
	return success
}

func (rc *RedisCache) Remove(key string) {
	rc.redisStore.Remove(key)
}

func (rc *RedisCache) IsExist(key string) bool {
	exist, err := rc.redisStore.IsExist(key)
	if err != nil {
		exist = false
	}
	return exist
}

func (rc *RedisCache) Len() int {
	return -1
}

// Cap return cache capacity
func (rc *RedisCache) Cap() int {
	return capUnlimit
}

// RealStore return the real redis store to do any other operation such as destroy
func (rc *RedisCache) RealStore() *redis.RedisStore {
	return rc.redisStore
}
