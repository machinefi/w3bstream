package job

import (
	"context"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type EventLogCleanupTask struct{}

func (t *EventLogCleanupTask) Arg() interface{} { return t }

func (t *EventLogCleanupTask) Subject() string { return "EventLogCleanup" }

func (t *EventLogCleanupTask) Handle(ctx context.Context) error {
	m := &models.EventLog{}
	d := types.MustMgrDBExecutorFromContext(ctx)

	_, err := d.Exec(builder.Delete().From(
		d.T(m),
		builder.Where(m.ColCreatedAt().Lt(time.Now().Add(-720*time.Hour))),
	))
	return err
}
