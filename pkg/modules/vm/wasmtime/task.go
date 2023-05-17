package wasmtime

import (
	"context"
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
}

func (t *Task) Handle(ctx context.Context) {
	panic("deprecated")
}
