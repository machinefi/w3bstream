package operator

import (
	"context"

	"github.com/ethereum/go-ethereum/crypto"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/projectoperator"
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

func GetDetailBySFID(ctx context.Context, id types.SFID) (*Detail, error) {
	o, err := GetBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	return convDetail(o)
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

func GetDetailByAccountAndName(ctx context.Context, accountID types.SFID, name string) (*Detail, error) {
	o, err := GetByAccountAndName(ctx, accountID, name)
	if err != nil {
		return nil, err
	}
	return convDetail(o)
}

func RemoveBySFID(ctx context.Context, id types.SFID) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Operator{RelOperator: models.RelOperator{OperatorID: id}}

	occupied, err := projectoperator.IsOperatorOccupied(ctx, id)
	if err != nil {
		return err
	}
	if occupied {
		return status.OccupiedOperator
	}

	if err := m.DeleteByOperatorID(d); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func Create(ctx context.Context, r *CreateReq) (*models.Operator, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	id := confid.MustSFIDGeneratorFromContext(ctx).MustGenSFID()
	acc := types.MustAccountFromContext(ctx)

	op := &models.Operator{
		RelAccount:  models.RelAccount{AccountID: acc.AccountID},
		RelOperator: models.RelOperator{OperatorID: id},
		OperatorInfo: models.OperatorInfo{
			Name:         r.Name,
			PrivateKey:   r.PrivateKey,
			PaymasterKey: r.PaymasterKey,
			Type:         r.Type,
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

func ListByCond(ctx context.Context, r *CondArgs) ([]models.Operator, error) {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.Operator{}
	)
	data, err := m.List(d, r.Condition())
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return data, nil
}

func List(ctx context.Context, r *ListReq) (*ListRsp, error) {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.Operator{}

		err  error
		ret  = &ListRsp{}
		cond = r.Condition()
		adds = r.Additions()
	)

	ret.Data, err = m.List(d, cond, adds...)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	ret.Total, err = m.Count(d, cond)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return ret, nil
}

func ListDetail(ctx context.Context, r *ListReq) (*ListDetailRsp, error) {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.Operator{}

		err  error
		ret  = &ListDetailRsp{}
		cond = r.Condition()
		adds = r.Additions()
	)

	data, err := m.List(d, cond, adds...)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	ret.Total, err = m.Count(d, cond)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}

	for i := range data {
		detail, err := convDetail(&data[i])
		if err != nil {
			return nil, err
		}
		ret.Data = append(ret.Data, *detail)
	}

	return ret, nil
}

func convDetail(d *models.Operator) (*Detail, error) {
	prvkey, err := crypto.HexToECDSA(d.PrivateKey)
	if err != nil {
		return nil, err
	}
	pubkey := crypto.PubkeyToAddress(prvkey.PublicKey)
	return &Detail{
		Operator: *d,
		Address:  pubkey.Hex(),
	}, nil
}
