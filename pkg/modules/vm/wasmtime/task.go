package wasmtime

import (
	"context"
	"fmt"
	"time"

	"github.com/iotexproject/Bumblebee/kit/mq"
	"github.com/iotexproject/w3bstream/pkg/types/wasm"
)

func NewTask(vm *Instance, fn string, pl []byte) *Task {
	return &Task{
		vm:      vm,
		Handler: fn,
		Payload: pl,
		Res:     make(chan *wasm.EventHandleResult, 1),
	}
}

type Task struct {
	vm      *Instance
	EventID string
	Handler string
	Payload []byte
	Res     chan *wasm.EventHandleResult
	mq.TaskState
}

var _ mq.Task = (*Task)(nil)

func (t *Task) Subject() string {
	return "HandleEvent"
}

func (t *Task) Arg() interface{} {
	return t
}

func (t *Task) Wait(du time.Duration) *wasm.EventHandleResult {
	select {
	case <-time.After(du):
		return &wasm.EventHandleResult{
			InstanceID: t.vm.ID(),
			Rsp:        nil,
			Code:       -1,
			ErrMsg:     "handle timeout",
		}
	case ret := <-t.Res:
		return ret
	}
}

func (t *Task) ID() string {
	return fmt.Sprintf("%s::%s::%s", t.Subject(), t.vm.ID(), t.EventID)
}

func (t *Task) Handle(ctx context.Context) {
	t.Res <- t.vm.Handle(ctx, t)
}
