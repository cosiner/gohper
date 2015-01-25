package server

import (
	"bytes"
	"encoding/gob"

	"github.com/cosiner/gomodule/redis"
)

type redisStore struct {
	store *redis.RedisStore
}

func (rstore *redisStore) Init(conf string) (err error) {
	if rstore.store == nil {
		rstore.store, err = redis.NewRedisStore(conf)
	}
	return
}

func (rstore *redisStore) IsExist(id string) bool {
	exist, _ := rstore.store.IsExist(id)
	return exist
}

func (rstore *redisStore) Save(id string, values map[string]interface{}, expire uint64) {
	var (
		buffer  = bytes.NewBuffer([]byte{})
		encoder = gob.NewEncoder(buffer)
	)
	if err := encoder.Encode(values); err == nil {
		go rstore.store.SetWithExpire(id, buffer.Bytes(), expire)
	}
}

func (rstore *redisStore) Get(id string) (vals map[string]interface{}) {
	if bs, err := redis.ToBytes(rstore.store.Get(id)); err == nil {
		vals = make(map[string]interface{})
		if len(bs) != 0 {
			decoder := gob.NewDecoder(bytes.NewBuffer(bs))
			err = decoder.Decode(vals)
		}
	}
	return
}

func (rstore *redisStore) Rename(oldId string, newId string) {
	rstore.store.Update("RENAME", oldId, newId)
}
