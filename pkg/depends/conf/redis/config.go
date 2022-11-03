package redis

import (
	"context"

	"github.com/gomodule/redigo/redis"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
)

type Config struct {
	Protocol       string
	Host           string
	Port           int
	Password       types.Password
	ConnectTimeout types.Duration
	WriteTimeout   types.Duration
	ReadTimeout    types.Duration
	IdleTimeout    types.Duration
	MaxActive      int
	MaxIdle        int
	Wait           bool
	DB             int
	pool           *redis.Pool
}

func (r *Config) Get() Conn { return nil }

func (r *Config) GetContext(ctx context.Context) (Conn, error) { return nil, nil }
