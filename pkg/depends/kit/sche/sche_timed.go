package sche

import (
	"context"
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sche/safe"
)

type timed struct {
	started *safe.Bool
	du      time.Duration
	fn      Job
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewTimedScheduler(fn Job, du time.Duration) Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &timed{
		started: safe.NewBool(),
		du:      du,
		fn:      fn,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func RunTimedScheduler(j Job, du time.Duration) Scheduler {
	ret := NewTimedScheduler(j, du)
	ret.Run()
	return ret
}

func (t *timed) Desc() string {
	return fmt.Sprintf("TIMED: %s", t.du)
}

func (t *timed) Context() context.Context { return t.ctx }

func (t *timed) WithContext(ctx context.Context) Scheduler {
	t.ctx, t.cancel = context.WithCancel(ctx)
	return t
}

func (t *timed) Started() bool { return t.started.Val() }

func (t *timed) Start() {
	if t.started.CAS(false, true) {
		go t.fn.Do()
		for {
			select {
			case <-t.ctx.Done():
				return
			case <-time.After(t.du):
				go t.fn.Do()
			}
		}
	}
}

func (t *timed) Run() { go t.Start() }

func (t *timed) Stop() {
	t.cancel()
}

func (t *timed) Value() interface{} {
	return t.fn
}
