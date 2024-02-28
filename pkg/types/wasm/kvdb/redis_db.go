package kvdb

import (
	"context"

	"github.com/gomodule/redigo/redis"

	confredis "github.com/machinefi/w3bstream/pkg/depends/conf/redis"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
)

type Cmd = confredis.Cmd

func NewRedisDB(d *confredis.Redis, prefix string) *RedisDB {
	return &RedisDB{
		d: d.WithPrefix(prefix),
	}
}

type RedisDB struct {
	d *confredis.Redis
}

func (r *RedisDB) Get(k string) ([]byte, error) {
	return redis.Bytes(r.d.Exec(&Cmd{Name: "GET", Args: []any{r.d.Key(k)}}))
}

func (r *RedisDB) Set(k string, v []byte) error {
	_, err := r.d.Exec(&Cmd{Name: "SET", Args: []any{r.d.Key(k), v}})
	return err
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
