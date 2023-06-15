package kvdb

import (
	"context"
	"strconv"

	"github.com/gomodule/redigo/redis"

	confredis "github.com/machinefi/w3bstream/pkg/depends/conf/redis"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
)

type RedisDB struct {
	db *confredis.Redis
}

func NewRedisDB(d *confredis.Redis) *RedisDB {
	return &RedisDB{db: d}
}

func (r *RedisDB) Get(key string) ([]byte, error) {
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

func (r *RedisDB) Set(key string, value []byte) error {
	var args []interface{}
	args = append(args, r.db.Prefix, key, string(value))
	if _, err := r.db.Exec(&confredis.Cmd{Name: "HSET", Args: args}); err != nil {
		return err
	}
	return nil
}

func (r *RedisDB) IncrBy(key string, value []byte) ([]byte, error) {
	var args []interface{}
	count, _ := strconv.Atoi(string(value))
	args = append(args, r.db.Key(key), count)
	result, err := r.db.Exec(&confredis.Cmd{Name: "INCRBY", Args: args})
	if err != nil || result == nil {
		return nil, err
	}
	val, err := redis.Int64(result, nil)
	if err != nil {
		return nil, err
	}
	return []byte(strconv.FormatInt(val, 10)), nil
}

// GetKey GET key
func (r *RedisDB) GetKey(key string) ([]byte, error) {
	var args []interface{}
	args = append(args, r.db.Key(key))
	result, err := r.db.Exec(&confredis.Cmd{Name: "GET", Args: args})
	if err != nil || result == nil {
		return nil, err
	}
	val, err := redis.Bytes(result, nil)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// SetKey SET key value
func (r *RedisDB) SetKey(key string, value []byte) error {
	var args []interface{}
	args = append(args, r.db.Key(key), string(value))
	if _, err := r.db.Exec(&confredis.Cmd{Name: "SET", Args: args}); err != nil {
		return err
	}
	return nil
}

func (r *RedisDB) SetKeyWithEX(key string, value []byte, exp int64) error {
	var args []interface{}
	args = append(args, r.db.Key(key), exp, string(value))
	if _, err := r.db.Exec(&confredis.Cmd{Name: "SETEX", Args: args}); err != nil {
		return err
	}
	return nil
}

type redisDBKey struct{}

func WithRedisDBKeyContext(redisDB *RedisDB) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, redisDBKey{}, redisDB)
	}
}

func RedisDBKeyFromContext(ctx context.Context) (*RedisDB, bool) {
	j, ok := ctx.Value(redisDBKey{}).(*RedisDB)
	return j, ok
}

func MustRedisDBKeyFromContext(ctx context.Context) *RedisDB {
	j, ok := ctx.Value(redisDBKey{}).(*RedisDB)
	must.BeTrue(ok)
	return j
}
