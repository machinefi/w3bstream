package tasks

import (
	"context"
	"reflect"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/modules/vm/wasmtime"
)

type HandleEvent struct {
	task *wasmtime.Task
}

func (t *HandleEvent) SetArg(v interface{}) error {
	if ctx, ok := v.(*wasmtime.Task); ok {
		t.task = ctx
		return nil
	}
	return errors.Errorf("invalid arg: %s", reflect.TypeOf(v))
}

func (t *HandleEvent) Output(ctx context.Context) (interface{}, error) {
	t.task.Handle(ctx)
	return nil, nil
}
