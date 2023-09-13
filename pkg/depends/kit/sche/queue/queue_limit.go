package queue

import (
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sche/safe"
)

type limited struct {
	qch    chan interface{}
	closed *safe.Bool
}

func NewLimited(lmt int) Queue {
	return &limited{
		qch:    make(chan interface{}, lmt),
		closed: safe.NewBool(),
	}
}

func (q *limited) Push(v interface{}) {
	q.qch <- v
}

func (q *limited) TryPush(v interface{}) bool {
	select {
	case q.qch <- v:
		return true
	default:
		return false
	}
}

func (q *limited) Pop() interface{} {
	return <-q.qch
}

func (q *limited) TryPop() interface{} {
	select {
	case ret := <-q.qch:
		return ret
	default:
		return nil
	}
}

func (q *limited) WaitPop(d time.Duration) interface{} {
	select {
	case <-time.After(d):
		return nil
	case ret := <-q.qch:
		return ret
	}
}

func (q *limited) Len() int {
	if q.closed.Val() {
		return 0
	}
	return len(q.qch)
}

func (q *limited) Close() {
	q.closed.Set(true)
	close(q.qch)
}

func (q *limited) Closed() bool { return q.closed.Val() }
