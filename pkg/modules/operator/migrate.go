package operator

import (
	"context"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"

	"github.com/pkg/errors"
)

// will delete at next version
func Migrate(ctx context.Context) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "operator.Migrate")
	defer l.End()

	a := &models.Account{}
	as, err := a.List(d, nil)
	if err != nil {
		l.Error(errors.Wrap(err, "list account failed"))
		return
	}
	for _, a := range as {
		if a.OperatorPrivateKey == "" {
			continue
		}

		id := confid.MustSFIDGeneratorFromContext(ctx).MustGenSFID()

		op := &models.Operator{
			RelAccount:  models.RelAccount{AccountID: a.AccountID},
			RelOperator: models.RelOperator{OperatorID: id},
			OperatorInfo: models.OperatorInfo{
				Name:       DefaultOperatorName,
				PrivateKey: a.OperatorPrivateKey,
			},
		}

		if err := op.Create(d); err != nil {
			if sqlx.DBErr(err).IsConflict() {
				continue
			}
			l.Error(errors.Wrap(err, "create operator failed"))
		}
	}
}
