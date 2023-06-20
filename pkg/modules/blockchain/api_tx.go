package blockchain

import (
	"context"
	"time"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
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

	if err := checkChainID(ctx, r.ChainID); err != nil {
		return nil, err
	}

	m := &models.ChainTx{
		RelChainTx: models.RelChainTx{ChainTxID: idg.MustGenSFID()},
		ChainTxData: models.ChainTxData{
			ProjectName: r.ProjectName,
			Uniq:        chainUniqFlag,
			Finished:    datatypes.FALSE,
			ChainTxInfo: r.ChainTxInfo,
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

func ListChainTxBySFIDs(ctx context.Context, ids []types.SFID) ([]models.ChainTx, error) {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	m := &models.ChainTx{}

	data, err := m.List(d, m.ColChainTxID().In(ids))
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return data, nil
}

func RemoveChainTxBySFID(ctx context.Context, id types.SFID) error {
	d := types.MustMonitorDBExecutorFromContext(ctx)

	m := &models.ChainTx{RelChainTx: models.RelChainTx{ChainTxID: id}}
	if err := m.DeleteByChainTxID(d); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func BatchUpdateChainTxPausedBySFIDs(ctx context.Context, ids []types.SFID, s datatypes.Bool) error {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	m := &models.ChainTx{
		ChainTxData: models.ChainTxData{
			ChainTxInfo: models.ChainTxInfo{
				Paused: s,
			},
		},
	}

	expr := builder.Update(d.T(m)).Set(
		m.ColPaused().ValueBy(s),
		m.ColUpdatedAt().ValueBy(types.Timestamp{Time: time.Now()}),
	).Where(m.ColChainTxID().In(ids))

	if _, err := d.Exec(expr); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}
