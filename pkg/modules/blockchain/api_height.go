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

type CreateChainHeightReq struct {
	ProjectName string `json:"-"`
	models.ChainHeightInfo
}

func CreateChainHeight(ctx context.Context, r *CreateChainHeightReq) (*models.ChainHeight, error) {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	if err := checkChainID(ctx, r.ChainID); err != nil {
		return nil, err
	}

	n := *r
	n.EventType = getEventType(n.EventType)
	m := &models.ChainHeight{
		RelChainHeight: models.RelChainHeight{ChainHeightID: idg.MustGenSFID()},
		ChainHeightData: models.ChainHeightData{
			ProjectName:     r.ProjectName,
			Uniq:            chainUniqFlag,
			Finished:        datatypes.FALSE,
			ChainHeightInfo: n.ChainHeightInfo,
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

func RemoveChainHeightBySFID(ctx context.Context, id types.SFID) error {
	d := types.MustMonitorDBExecutorFromContext(ctx)

	m := &models.ChainHeight{RelChainHeight: models.RelChainHeight{ChainHeightID: id}}
	if err := m.DeleteByChainHeightID(d); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}
