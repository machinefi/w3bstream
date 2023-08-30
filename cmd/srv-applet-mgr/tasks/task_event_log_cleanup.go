package tasks

import (
	"context"
	"reflect"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/modules/job"
)

type EventLogCleanup struct {
	*job.EventLogCleanupTask
}

func (t *EventLogCleanup) SetArg(v interface{}) error {
	if ctx, ok := v.(*job.EventLogCleanupTask); ok {
		t.EventLogCleanupTask = ctx
		return nil
	}
	return errors.Errorf("invalid arg: %s", reflect.TypeOf(v))
}

func (t *EventLogCleanup) Output(ctx context.Context) (interface{}, error) {
	ctx, l := logr.Start(ctx, "tasks.EventLog.Output")
	defer l.End()

	return nil, t.Handle(ctx)
}
