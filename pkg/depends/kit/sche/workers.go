package sche

import (
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sche/queue"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sche/safe"
)

const defaultPoolLimit = 64

type Workers struct {
	q   queue.Queue
	lmt int
	seq *safe.Int64
}

func NewWorkers(lmt ...int) *Workers {
	limit := defaultPoolLimit
	if len(lmt) > 0 && lmt[0] > 0 {
		limit = lmt[0]
	}
	return &Workers{q: queue.NewLimited(limit), lmt: limit, seq: safe.NewInt64()}
}

func (p *Workers) Limit() int { return p.lmt }

func (p *Workers) Add(j Job) (ctx *Context) {
	ctx = NewContext(j, p.seq.Add(1))
	p.q.Push(ctx)
	return ctx
}

func (p *Workers) AddWithDeadline(j Job, deadline time.Time) (ctx *Context) {
	ctx = NewContext(j, p.seq.Add(1))
	p.q.Push(ctx.WithDeadline(deadline))
	return ctx
}

func (p *Workers) AddWithTimeout(j Job, timeout time.Duration) (ctx *Context) {
	ctx = NewContext(j, p.seq.Add(1))
	p.q.Push(ctx.WithTimeout(timeout))
	return ctx
}

func (p *Workers) Pop() *Context {
	ctx := p.q.Pop()
	if ctx == nil {
		return nil
	}
	ctx.(*Context).stat[1] = time.Now()
	return ctx.(*Context)
}

func (p *Workers) Close() { p.q.Close() }

func (p *Workers) Closed() bool { return p.q.Closed() }

func (p *Workers) Len() int { return p.q.Len() }
