package tasks

import (
	"context"
	"reflect"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/modules/job"
)

type EventLog struct {
	*job.EventLogTask
}

func (t *EventLog) SetArg(v interface{}) error {
	if ctx, ok := v.(*job.EventLogTask); ok {
		t.EventLogTask = ctx
		return nil
	}
	return errors.Errorf("invalid arg: %s", reflect.TypeOf(v))
}

func (t *EventLog) Output(ctx context.Context) (interface{}, error) {
	ctx, l := logr.Start(ctx, "tasks.EventLog.Output")
	defer l.End()

	return nil, t.Handle(ctx)
}
