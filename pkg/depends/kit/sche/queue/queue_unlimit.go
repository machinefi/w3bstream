package queue

import (
	"math"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sche/safe"
)

type unlimited struct {
	lst    *safe.List
	act    chan struct{}
	qch    chan interface{}
	closed *safe.Bool
}

func New() Queue {
	ret := &unlimited{
		lst:    safe.NewList(),
		act:    make(chan struct{}, math.MaxInt32),
		qch:    make(chan interface{}, gCap),
		closed: safe.NewBool(),
	}
	go ret.sync()
	return ret
}

const (
	gCap   = 4096
	gBatch = 32
)

func (q *unlimited) Push(v interface{}) {
	if !q.closed.Val() {
		q.lst.PushBack(v)
		q.act <- struct{}{}
	}
}

func (q *unlimited) TryPush(v interface{}) bool { q.Push(v); return true }

func (q *unlimited) Pop() interface{} { return <-q.qch }

func (q *unlimited) TryPop() interface{} {
	select {
	case ret := <-q.qch:
		return ret
	default:
		return nil
	}
}

func (q *unlimited) WaitPop(d time.Duration) interface{} {
	select {
	case <-time.After(d):
		return nil
	case ret := <-q.qch:
		return ret
	}
}

func (q *unlimited) Len() int {
	return len(q.qch) + len(q.act)
}

func (q *unlimited) Close() {
	q.closed.Set(true)
	close(q.act)
	close(q.qch)
	q.lst.Clear()
}

func (q *unlimited) sync() {
	defer func() {
		if q.closed.Val() {
			_ = recover()
		}
	}()
	for !q.closed.Val() {
		<-q.act
		if !q.closed.Val() {
			bat := q.lst.PopFrontN(gBatch)
			for _, v := range bat {
				q.qch <- v
			}
			for i := 0; i < len(bat); i++ {
				<-q.act
			}
		} else {
			break
		}
	}
}

func (q *unlimited) Closed() bool { return q.closed.Val() }
