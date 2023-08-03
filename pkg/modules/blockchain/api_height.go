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

type CreateChainHeightReq struct {
	ProjectName string `json:"-"`
	models.ChainHeightInfo
}

func CreateChainHeight(ctx context.Context, r *CreateChainHeightReq) (*models.ChainHeight, error) {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	if err := checkChain(ctx, r.ChainID, r.ChainName); err != nil {
		return nil, err
	}

	m := &models.ChainHeight{
		RelChainHeight: models.RelChainHeight{ChainHeightID: idg.MustGenSFID()},
		ChainHeightData: models.ChainHeightData{
			ProjectName:     r.ProjectName,
			Uniq:            chainUniqFlag,
			Finished:        datatypes.FALSE,
			ChainHeightInfo: r.ChainHeightInfo,
		},
	}
	if err := m.Create(d); err != nil {
		if sqlx.DBErr(err).IsConflict() {
			return nil, status.ChainHeightConflict
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func GetChainHeightBySFID(ctx context.Context, id types.SFID) (*models.ChainHeight, error) {
	d := types.MustMonitorDBExecutorFromContext(ctx)

	m := &models.ChainHeight{RelChainHeight: models.RelChainHeight{ChainHeightID: id}}
	if err := m.FetchByChainHeightID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.ChainHeightNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func ListChainHeightBySFIDs(ctx context.Context, ids []types.SFID) ([]models.ChainHeight, error) {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	m := &models.ChainHeight{}

	data, err := m.List(d, m.ColChainHeightID().In(ids))
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return data, nil
}

func RemoveChainHeightBySFID(ctx context.Context, id types.SFID) error {
	d := types.MustMonitorDBExecutorFromContext(ctx)

	m := &models.ChainHeight{RelChainHeight: models.RelChainHeight{ChainHeightID: id}}
	if err := m.DeleteByChainHeightID(d); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func BatchUpdateChainHeightPausedBySFIDs(ctx context.Context, ids []types.SFID, s datatypes.Bool) error {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	m := &models.ChainHeight{
		ChainHeightData: models.ChainHeightData{
			ChainHeightInfo: models.ChainHeightInfo{
				Paused: s,
			},
		},
	}

	expr := builder.Update(d.T(m)).Set(
		m.ColPaused().ValueBy(s),
		m.ColUpdatedAt().ValueBy(types.Timestamp{Time: time.Now()}),
	).Where(m.ColChainHeightID().In(ids))

	if _, err := d.Exec(expr); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}
