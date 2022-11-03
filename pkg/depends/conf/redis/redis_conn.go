package redis

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

type Conn = redis.Conn

func Command(name string, args ...interface{}) *Cmd {
	return &Cmd{
		Name: name,
		Args: args,
	}
}

type Cmd struct {
	Name string
	Args []interface{}
}

type Operator interface {
	Prefix(key string) string
	Get() Conn
	GetContext(ctx context.Context) (Conn, error)
	Exec(cmd *Cmd, others ...*Cmd) (interface{}, error)
	ExecContext(ctx context.Context, cmd *Cmd, others ...*Cmd) (interface{}, error)
}
