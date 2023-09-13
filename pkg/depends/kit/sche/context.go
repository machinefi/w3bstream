package sche

import (
	"context"
	"time"
)

type Context struct {
	stat                   // stat commit,committed,execute,finished
	deadline time.Time     // deadline context deadline
	result   chan *Result  // result chan
	res      *Result       // res result
	done     chan struct{} // done context done chan
	seq      int64         // seq job sequence in workers pool
	Job
}

func NewContext(j Job, seq int64) *Context {
	ret := &Context{
		Job:    j,
		result: make(chan *Result, 1),
		res:    &Result{},
		done:   make(chan struct{}, 1),
		seq:    seq,
	}
	ret.stat[0] = time.Now()
	return ret
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	if !c.deadline.IsZero() {
		return c.deadline, true
	}
	return
}

func (c *Context) Value(_ interface{}) interface{} { return nil }

func (c *Context) Done() <-chan struct{} { return c.done }

func (c *Context) Err() error {
	if c.res == nil {
		return nil
	}
	return c.res.error
}

func (c *Context) Exec(ctx context.Context) {
	if !c.deadline.IsZero() {
		ctx, _ = context.WithDeadline(ctx, c.deadline)
	}
	select {
	case <-ctx.Done():
		c.res.error = ctx.Err()
	default:
		c.res = &Result{}
		c.res.Val, c.res.error = c.Job.Do()
		c.stat[2] = time.Now()
		c.result <- c.res
	}
	c.done <- struct{}{}
}

func (c *Context) WithDeadline(deadline time.Time) *Context {
	c.deadline = deadline
	return c
}

func (c *Context) WithTimeout(timeout time.Duration) *Context {
	c.deadline = c.stat[0].Add(timeout)
	return c
}

func (c *Context) Result() (interface{}, error) {
	r := <-c.result
	return r.Val, r.error
}

func (c *Context) Sequence() int64 { return c.seq }

// stat 0 commit 1 scheduled 3 done
type stat [3]time.Time

// Latency sub committed and commit
func (s *stat) Latency() time.Duration {
	if !s[0].IsZero() && !s[1].IsZero() {
		return s[1].Sub(s[0])
	}
	return -1
}

// Cost sub done and scheduled
func (s *stat) Cost() time.Duration {
	if !s[1].IsZero() && !s[2].IsZero() {
		return s[2].Sub(s[1])

	}
	return -1
}

// Total sub commit and
func (s *stat) Total() time.Duration {
	if !s[0].IsZero() && !s[2].IsZero() {
		return s[2].Sub(s[0])

	}
	return -1
}

func (s *stat) Stat() [2]int64 {
	return [2]int64{s.Latency().Milliseconds(), s.Cost().Milliseconds()}
}
