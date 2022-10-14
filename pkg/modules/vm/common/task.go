package common

import "github.com/iotexproject/w3bstream/pkg/types/wasm"

type Task struct {
	Handler string
	Payload []byte
	Res     chan *EventHandleResult
}

type EventHandleResult struct {
	Response []byte
	Code     wasm.ResultStatusCode
}

type TaskQueue struct{ Ch chan *Task }

func (t *TaskQueue) Wait() <-chan *Task { return t.Ch }

func (t *TaskQueue) Push(task *Task) { t.Ch <- task }
