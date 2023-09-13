package sche

import (
	"context"
	"sync"
	"time"
)

type Fn func()

func (f Fn) Do() (interface{}, error) { f(); return nil, nil }

type FnWithErr func() error

func (f FnWithErr) Do() (interface{}, error) { return nil, f() }

type FnWithVal func() interface{}

func (f FnWithVal) Do() (interface{}, error) { return f(), nil }

type FnWithRes func() (interface{}, error)

func (f FnWithRes) Do() (interface{}, error) { return f() }

// Job sche schedule unit
type Job interface{ Do() (interface{}, error) }

// Result job result
type Result struct {
	Val interface{}
	error
}

// Scheduler scheduler
type Scheduler interface {
	Start()
	Run()
	Started() bool
	Stop()
	Context() context.Context
	WithContext(ctx context.Context) Scheduler
}

// WorkersScheduler scheduler with worker
type WorkersScheduler interface {
	Scheduler
	Add(Job) *Context
	AddWithDeadline(Job, time.Time) *Context
	AddWithTimeout(Job, time.Duration) *Context
	WaitGroup(...Job) []*Context
}

// WaitGroup for batch commit jobs
func WaitGroup(sche WorkersScheduler, jobs ...Job) (ret []*Context) {
	ch := make(chan *Context, len(jobs))
	wg := &sync.WaitGroup{}
	for i := range jobs {
		select {
		case <-sche.Context().Done():
			return nil
		default:
			Group(wg).Do(func() {
				ctx := sche.Add(jobs[i])
				ch <- ctx
				<-ctx.Done()
			})
		}
	}
	wg.Wait()
	for i := 0; i < len(jobs); i++ {
		ret = append(ret, <-ch)
	}
	return
}
