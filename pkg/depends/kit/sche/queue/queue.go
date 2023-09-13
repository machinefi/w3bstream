package queue

import "time"

type Queue interface {
	Push(v interface{})
	TryPush(v interface{}) bool
	Pop() interface{}
	TryPop() interface{}
	WaitPop(time.Duration) interface{}
	Len() int
	Close()
	Closed() bool
}
