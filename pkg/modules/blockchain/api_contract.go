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

type CreateContractLogReq struct {
	ProjectName string `json:"-"`
	models.ContractLogInfo
}

func CreateContractLog(ctx context.Context, r *CreateContractLogReq) (*models.ContractLog, error) {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	if err := checkChainID(ctx, r.ChainID); err != nil {
		return nil, err
	}

	r.BlockCurrent = r.BlockStart
	m := &models.ContractLog{
		RelContractLog: models.RelContractLog{ContractLogID: idg.MustGenSFID()},
		ContractLogData: models.ContractLogData{
			ProjectName:     r.ProjectName,
			Uniq:            chainUniqFlag,
			ContractLogInfo: r.ContractLogInfo,
		},
	}
	if err := m.Create(d); err != nil {
		if sqlx.DBErr(err).IsConflict() {
			return nil, status.ContractLogConflict
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func checkChainID(ctx context.Context, id uint64) error {
	ethcli := types.MustETHClientConfigFromContext(ctx)
	if _, ok := ethcli.Clients[uint32(id)]; !ok {
		return status.BlockchainNotFound
	}
	return nil
}

func GetContractLogBySFID(ctx context.Context, id types.SFID) (*models.ContractLog, error) {
	d := types.MustMonitorDBExecutorFromContext(ctx)

	m := &models.ContractLog{RelContractLog: models.RelContractLog{ContractLogID: id}}
	if err := m.FetchByContractLogID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.ContractLogNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func ListContractLogBySFIDs(ctx context.Context, ids []types.SFID) ([]models.ContractLog, error) {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	m := &models.ContractLog{}

	data, err := m.List(d, m.ColContractLogID().In(ids))
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return data, nil
}

func RemoveContractLogBySFID(ctx context.Context, id types.SFID) error {
	d := types.MustMonitorDBExecutorFromContext(ctx)

	m := &models.ContractLog{RelContractLog: models.RelContractLog{ContractLogID: id}}
	if err := m.DeleteByContractLogID(d); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func BatchUpdateContractLogPausedBySFIDs(ctx context.Context, ids []types.SFID, s datatypes.Bool) error {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	m := &models.ContractLog{
		ContractLogData: models.ContractLogData{
			ContractLogInfo: models.ContractLogInfo{
				Paused: s,
			},
		},
	}

	expr := builder.Update(d.T(m)).Set(
		m.ColPaused().ValueBy(s),
		m.ColUpdatedAt().ValueBy(types.Timestamp{Time: time.Now()}),
	).Where(m.ColContractLogID().In(ids))

	if _, err := d.Exec(expr); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}
