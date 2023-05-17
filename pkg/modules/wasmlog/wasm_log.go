package wasmlog

import (
	"context"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func RemoveBySFID(ctx context.Context, id types.SFID) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.WasmLog{RelWasmLog: models.RelWasmLog{WasmLogID: id}}

	if err := m.DeleteByWasmLogID(d); err != nil {
		return status.DatabaseError.StatusErr().
			WithDesc(errors.Wrap(err, id.String()).Error())
	}
	return nil
}

func Remove(ctx context.Context, r *CondArgs) error {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.WasmLog{}

		err error
		lst []models.WasmLog
	)

	return sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			lst, err = m.List(d, r.Condition())
			if err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			ctx := types.WithMgrDBExecutor(ctx, d)
			summary := statusx.ErrorFields{}
			for i := range lst {
				v := &lst[i]
				if err = RemoveBySFID(ctx, v.WasmLogID); err != nil {
					se := statusx.FromErr(err)
					summary = append(summary, &statusx.ErrorField{
						In:    v.WasmLogID.String(),
						Field: se.Key,
						Msg:   se.Desc,
					})
				}
			}
			if len(summary) > 0 {
				return status.BatchRemoveWasmLogFailed.StatusErr().
					AppendErrorFields(summary...)
			}
			return nil
		},
	).Do()
}
