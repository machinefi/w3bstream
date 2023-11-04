package sche

import (
	"context"
	"fmt"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sche/safe"
)

type chain struct {
	*Workers
	id      string
	started *safe.Bool
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewChainedScheduler(lmt ...int) WorkersScheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &chain{
		Workers: NewWorkers(lmt...),
		started: safe.NewBool(),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func RunChainedScheduler(lmt ...int) WorkersScheduler {
	ret := NewChainedScheduler(lmt...)
	ret.Run()
	return ret
}

func (c *chain) Desc() string {
	return fmt.Sprintf("CHAINED: %d", c.Workers.Limit())
}

func (c *chain) Context() context.Context { return c.ctx }

func (c *chain) WithContext(ctx context.Context) Scheduler {
	c.ctx, c.cancel = context.WithCancel(ctx)
	return c
}

func (c *chain) Start() {
	if c.started.CAS(false, true) {
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				if ctx := c.Pop(); ctx != nil {
					ctx.Exec(c.ctx)
				}
			}
		}
	}
}

func (c *chain) Run() { go c.Start() }

func (c *chain) Started() bool { return c.started.Val() }

func (c *chain) Stop() {
	c.cancel()
	c.Workers.Close()
}

func (c *chain) WaitGroup(jobs ...Job) []*Context {
	return WaitGroup(c, jobs...)
}
