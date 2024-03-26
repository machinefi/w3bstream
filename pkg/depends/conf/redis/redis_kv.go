package redis

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
)

type CacheKeys interface {
	CacheKeys() []string
}

func (r *Redis) Get(k string, v any) error {
	raw, err := redis.Bytes(r.Exec(&Cmd{Name: "GET", Args: []any{r.Key(k)}}))
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, v)
}

func (r *Redis) SetEx(k string, v any, exp time.Duration) error {
	raw, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = r.Exec(&Cmd{Name: "SET", Args: []any{r.Key(k), raw, "PX", exp.Milliseconds()}})
	return err
}

func (r *Redis) Set(k string, v any) error {
	raw, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = r.Exec(&Cmd{Name: "SET", Args: []any{r.Key(k), raw}})
	return err
}

func (r *Redis) Del(ks ...string) error {
	keys := make([]any, 0, len(ks))
	for _, k := range ks {
		keys = append(keys, r.Key(k))
	}
	_, err := r.Exec(&Cmd{Name: "DEL", Args: keys})
	return err
}

func (r *Redis) RawDel(ks ...string) error {
	keys := make([]any, 0, len(ks))
	for _, k := range ks {
		keys = append(keys, k)
	}
	_, err := r.Exec(&Cmd{Name: "DEL", Args: keys})
	return err
}

func (r *Redis) IncrBy(k string, count int64) (int64, error) {
	return redis.Int64(r.Exec(&Cmd{Name: "INCRBY", Args: []any{r.Key(k), count}}))
}

func (r *Redis) Exists(k string) (bool, error) {
	return redis.Bool(r.Exec(&Cmd{Name: "EXISTS", Args: []any{r.Key(k)}}))
}

func (r *Redis) Keys(pattern string) ([]string, error) {
	return redis.Strings(r.Exec(&Cmd{Name: "KEYS", Args: []any{r.Key(pattern)}}))
}
