package wasmtime

import (
	"context"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

type Task struct {
	EventID   string
	EventType string
	Handler   string
	Payload   []byte
	mq.TaskState

	vm       *Instance
	retrieve chan *wasm.EventHandleResult
}

var _ mq.Task = (*Task)(nil)

func (t *Task) Subject() string { return "HandleEvent" }

func (t *Task) ID() string { return t.EventID }

func (t *Task) Arg() interface{} { return t }

func (t *Task) Handle(ctx context.Context) {
	t.retrieve <- t.vm.handle(ctx, t)
}

func (t *Task) Wait() *wasm.EventHandleResult {
	select {
	case v := <-t.retrieve:
		return v
	case <-time.After(5 * time.Second):
		return &wasm.EventHandleResult{
			InstanceID: t.vm.ID(),
			Code:       wasm.ResultStatusCode_Failed,
			ErrMsg:     "wait timeout",
		}
	}
}
