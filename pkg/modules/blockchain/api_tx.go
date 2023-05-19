package blockchain

import (
	"context"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateChainTxReq struct {
	ProjectName string `json:"-"`
	models.ChainTxInfo
}

func CreateChainTx(ctx context.Context, r *CreateChainTxReq) (*models.ChainTx, error) {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	if err := checkChainID(d, r.ChainID); err != nil {
		return nil, err
	}

	n := *r
	n.EventType = getEventType(n.EventType)
	m := &models.ChainTx{
		RelChainTx: models.RelChainTx{ChainTxID: idg.MustGenSFID()},
		ChainTxData: models.ChainTxData{
			ProjectName: r.ProjectName,
			Uniq:        chainUniqFlag,
			Finished:    datatypes.FALSE,
			ChainTxInfo: n.ChainTxInfo,
		},
	}
	if err := m.Create(d); err != nil {
		if sqlx.DBErr(err).IsConflict() {
			return nil, status.ChainTxConflict
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func GetChainTxBySFID(ctx context.Context, id types.SFID) (*models.ChainTx, error) {
	d := types.MustMonitorDBExecutorFromContext(ctx)

	m := &models.ChainTx{RelChainTx: models.RelChainTx{ChainTxID: id}}
	if err := m.FetchByChainTxID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.ChainTxNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func RemoveChainTxBySFID(ctx context.Context, id types.SFID) error {
	d := types.MustMonitorDBExecutorFromContext(ctx)

	m := &models.ChainTx{RelChainTx: models.RelChainTx{ChainTxID: id}}
	if err := m.DeleteByChainTxID(d); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}
