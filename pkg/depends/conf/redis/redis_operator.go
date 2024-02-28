package redis

import "github.com/gomodule/redigo/redis"

type Operator interface {
	// Key returns key with prefix
	Key(key string) string
	// Acquire redis connection
	Acquire() redis.Conn
	// Exec to execute redis command
	Exec(cmd *Cmd, others ...*Cmd) (interface{}, error)
}

var (
	_ Operator = (*Redis)(nil)
	_ Operator = (*Endpoint)(nil)
)
