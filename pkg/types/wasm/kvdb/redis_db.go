package kvdb

import (
	"github.com/gomodule/redigo/redis"

	confredis "github.com/machinefi/w3bstream/pkg/depends/conf/redis"
)

type redisDB struct {
	db *confredis.Redis
}

func NewRedisDB(d *confredis.Redis) *redisDB {
	return &redisDB{db: d}
}

func (r *redisDB) Get(key string) ([]byte, error) {
	var args []interface{}
	args = append(args, r.db.Prefix, key)
	result, err := r.db.Exec(&confredis.Cmd{Name: "HGET", Args: args})
	if err != nil || result == nil {
		return nil, err
	}
	val, err := redis.Bytes(result, nil)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (r *redisDB) Set(key string, value []byte) error {
	var args []interface{}
	args = append(args, r.db.Prefix, key, string(value))
	if _, err := r.db.Exec(&confredis.Cmd{Name: "HSET", Args: args}); err != nil {
		return err
	}
	return nil
}
