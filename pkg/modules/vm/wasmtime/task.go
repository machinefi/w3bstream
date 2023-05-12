package wasmtime

import (
	"context"
	"time"

	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func newTask(ctx context.Context, fn string, eventType string, data []byte) *Task {
	return &Task{
		ctx:       ctx,
		EventType: eventType,
		Handler:   fn,
		Payload:   data,
	}
}

type Task struct {
	ctx       context.Context
	EventID   string
	EventType string
	Handler   string
	Payload   []byte
	// mq.TaskState
}

// var _ mq.Task = (*Task)(nil)

func (t *Task) Subject() string {
	return "HandleEvent"
}

func (t *Task) Arg() interface{} {
	return t
}

func (t *Task) Wait(du time.Duration) *wasm.EventHandleResult {
	panic("deprecated")
	// select {
	// case <-time.After(du):
	// 	return &wasm.EventHandleResult{
	// 		InstanceID: t.vm.ID(),
	// 		Rsp:        nil,
	// 		Code:       -1,
	// 		ErrMsg:     "handle timeout",
	// 	}
	// case ret := <-t.Res:
	// 	return ret
	// }
}

func (t *Task) ID() string {
	panic("deprecated")
	// return fmt.Sprintf("%s::%s::%s", t.Subject(), t.vm.ID(), t.EventID)
}

func (t *Task) Handle(ctx context.Context) {
	// t.Res <- t.vm.Handle(ctx, t)
}
