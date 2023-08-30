package job

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func NewEventLogTask(m *models.EventLog) *EventLogTask {
	return &EventLogTask{
		EventLog:  m,
		TaskState: mq.TASK_STATE__PENDING,
	}
}

type EventLogTask struct {
	*models.EventLog
	mq.TaskState
	mq.TaskUUID
}

func (t *EventLogTask) Arg() interface{} { return t }

func (t *EventLogTask) Subject() string { return "EventLog" }

func (t *EventLogTask) Handle(ctx context.Context) error {
	if t.EventLog == nil {
		return nil
	}
	return t.EventLog.Create(types.MustMgrDBExecutorFromContext(ctx))
}
