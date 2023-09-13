package sche

import (
	"context"
	"fmt"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sche/safe"
)

type concurrent struct {
	*Workers
	started *safe.Bool
	con     int
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewConScheduler(concurrency int) WorkersScheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &concurrent{
		Workers: NewWorkers(concurrency),
		started: safe.NewBool(),
		con:     concurrency,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func RunConScheduler(concurrency int) WorkersScheduler {
	ret := NewConScheduler(concurrency)
	ret.Run()
	return ret
}

func (c *concurrent) Desc() string {
	return fmt.Sprintf("CONC: %d", c.con)
}

func (c *concurrent) Context() context.Context { return c.ctx }

func (c *concurrent) WithContext(ctx context.Context) Scheduler {
	c.ctx, c.cancel = context.WithCancel(ctx)
	return c
}

func (c *concurrent) routine() {
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

func (c *concurrent) run() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			if ctx := c.Pop(); ctx != nil {
				go ctx.Exec(c.ctx)
			}
		}
	}
}

func (c *concurrent) Started() bool { return c.started.Val() }

func (c *concurrent) Start() {
	if c.started.CAS(false, true) {
		if c.con == 0 {
			c.run()
		} else {
			workers := make([]func(), 0, c.con)
			for i := 0; i < c.con; i++ {
				workers = append(workers, c.routine)
			}
			GroupDoMany(workers...)
		}
	}
}

func (c *concurrent) Run() { go c.Start() }

func (c *concurrent) Stop() {
	c.cancel()
	c.Workers.Close()
}

func (c *concurrent) WaitGroup(jobs ...Job) []*Context {
	return WaitGroup(c, jobs...)
}
