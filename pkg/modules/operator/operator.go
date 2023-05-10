package operator

import (
	"context"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

const DefaultOperatorName = "default"

func GetBySFID(ctx context.Context, id types.SFID) (*models.Operator, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Operator{RelOperator: models.RelOperator{OperatorID: id}}

	if err := m.FetchByOperatorID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.OperatorNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func GetByAccountAndName(ctx context.Context, accountID types.SFID, name string) (*models.Operator, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Operator{
		OperatorInfo: models.OperatorInfo{Name: name},
		RelAccount:   models.RelAccount{AccountID: accountID},
	}

	if err := m.FetchByAccountIDAndName(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.OperatorNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func RemoveBySFID(ctx context.Context, id types.SFID) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Operator{RelOperator: models.RelOperator{OperatorID: id}}

	if err := m.DeleteByOperatorID(d); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func Create(ctx context.Context, r *CreateReq) (*models.Operator, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	id := confid.MustSFIDGeneratorFromContext(ctx).MustGenSFID()

	op := &models.Operator{
		RelAccount:  models.RelAccount{AccountID: r.AccountID},
		RelOperator: models.RelOperator{OperatorID: id},
		OperatorInfo: models.OperatorInfo{
			Name:       r.Name,
			PrivateKey: r.PrivateKey,
		},
	}

	if err := op.Create(d); err != nil {
		if sqlx.DBErr(err).IsConflict() {
			return nil, status.OperatorConflict
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return op, nil
}

func List(ctx context.Context, r *ListReq) (*ListRsp, error) {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.Operator{}

		err  error
		ret  = &ListRsp{}
		cond = r.Condition()
	)

	ret.Data, err = m.List(d, cond)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	ret.Total, err = m.Count(d, cond)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return ret, nil
}
