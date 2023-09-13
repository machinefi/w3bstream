package mq

import (
	"context"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq/mem_mq"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
)

type Config struct {
	// Channel worker channel name
	Channel string
	// Store type in memory
	Store StoreType
	// Limit message queue max length
	Limit int
	// WorkerCount task worker count (concurrency)
	WorkerCount      int
	PushQueueTimeout types.Duration

	*mq.TaskBoard
	*mq.TaskWorker
}

func (c *Config) SetDefault() {
	if c.Store == STORE_TYPE_UNKNOWN || c.Store > STORE_TYPE__REDIS {
		c.Store = STORE_TYPE__MEM
	}
	if c.PushQueueTimeout == 0 {
		c.PushQueueTimeout = types.Duration(time.Second)
	}
	if c.Limit == 0 {
		c.Limit = 1024
	}
	if c.WorkerCount == 0 {
		c.WorkerCount = 256
	}
}

func (c *Config) Init() error {
	if c.Channel == "" {
		c.Channel = os.Getenv(consts.EnvProjectName)
	}

	switch c.Store {
	case STORE_TYPE__MEM:
		tm := mem_mq.New(c.Limit)
		c.TaskWorker = mq.NewTaskWorker(tm, mq.WithWorkerCount(c.WorkerCount), mq.WithChannel(c.Channel))
		c.TaskBoard = mq.NewTaskBoard(tm)
		return nil
	default:
		return errors.New("mq Config init failed")
	}
}

func (c *Config) Name() string { return "TaskBoard" }

type kvc struct{}

func WithMq(ctx context.Context, c *Config) context.Context {
	return contextx.WithValue(ctx, kvc{}, c)
}

func WithMqContext(c *Config) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return WithMq(ctx, c)
	}
}

func MqFromContext(ctx context.Context) (*Config, bool) {
	v, ok := ctx.Value(kvc{}).(*Config)
	return v, ok
}

func MustMqFromContext(ctx context.Context) *Config {
	v, ok := MqFromContext(ctx)
	must.BeTrue(ok)
	return v
}
