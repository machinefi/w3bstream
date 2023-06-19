package blockchain

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

const chainUniqFlag = 0

type UpdateMonitorReq struct {
	IDs []types.SFID `json:"ids"`
}

func RemoveMonitor(ctx context.Context, projectName string) error {
	d := types.MustMonitorDBExecutorFromContext(ctx)

	return sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			m := &models.ContractLog{}

			expr := builder.Delete().From(d.T(m), builder.Where(m.ColProjectName().Eq(projectName)))
			if _, err := d.Exec(expr); err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			m := &models.ChainTx{}

			expr := builder.Delete().From(d.T(m), builder.Where(m.ColProjectName().Eq(projectName)))
			if _, err := d.Exec(expr); err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			m := &models.ChainHeight{}

			expr := builder.Delete().From(d.T(m), builder.Where(m.ColProjectName().Eq(projectName)))
			if _, err := d.Exec(expr); err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
	).Do()
}

func getEventType(eventType string) string {
	if eventType == "" {
		return enums.MONITOR_EVENTTYPEDEFAULT
	}
	return eventType
}

func getPaused(i datatypes.Bool) datatypes.Bool {
	if i == datatypes.TRUE {
		return datatypes.TRUE
	}
	return datatypes.FALSE
}
