package project

import (
	"context"

	confid "github.com/iotexproject/Bumblebee/conf/id"
	"github.com/iotexproject/Bumblebee/kit/sqlx"

	"github.com/iotexproject/w3bstream/pkg/errors/status"
	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type CreateMonitorReq struct {
	Contractlog *CreateContractlogReq `json:"contractLog,omitempty"`
	Chaintx     *CreateChaintxReq     `json:"chainTx,omitempty"`
	ChainHeight *CreateChainHeightReq `json:"chainHeight,omitempty"`
}

type (
	CreateContractlogReq = models.ContractlogInfo
	CreateChaintxReq     = models.ChaintxInfo
	CreateChainHeightReq = models.ChainHeightInfo
)

func CreateMonitor(ctx context.Context, projectName string, r *CreateMonitorReq) (interface{}, error) {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)
	switch {
	case r.Contractlog != nil:
		return createContractLog(d, projectName, r.Contractlog, idg)
	case r.Chaintx != nil:
		return createChainTx(d, projectName, r.Chaintx, idg)
	case r.ChainHeight != nil:
		return createChainHeight(d, projectName, r.ChainHeight, idg)
	default:
		return nil, status.BadRequest
	}

}

func createContractLog(d sqlx.DBExecutor, projectName string, r *CreateContractlogReq, idg confid.SFIDGenerator) (*models.Contractlog, error) {
	if err := checkChainID(d, r.ChainID); err != nil {
		return nil, err
	}

	n := *r
	n.BlockCurrent = n.BlockStart
	n.EventType = types.EVENTTYPEDEFAULT // TODO support event type
	m := &models.Contractlog{
		RelContractlog: models.RelContractlog{ContractlogID: idg.MustGenSFID()},
		ContractlogData: models.ContractlogData{
			ProjectName:     projectName,
			ContractlogInfo: n,
		},
	}
	if err := m.Create(d); err != nil {
		return nil, status.CheckDatabaseError(err, "CreateContractlogMonitor")
	}
	return m, nil
}

func createChainTx(d sqlx.DBExecutor, projectName string, r *CreateChaintxReq, idg confid.SFIDGenerator) (*models.Chaintx, error) {
	if err := checkChainID(d, r.ChainID); err != nil {
		return nil, err
	}

	n := *r
	n.EventType = types.EVENTTYPEDEFAULT // TODO support event type
	m := &models.Chaintx{
		RelChaintx: models.RelChaintx{ChaintxID: idg.MustGenSFID()},
		ChaintxData: models.ChaintxData{
			ProjectName: projectName,
			ChaintxInfo: n,
		},
	}
	if err := m.Create(d); err != nil {
		return nil, status.CheckDatabaseError(err, "CreateChainTxMonitor")
	}
	return m, nil
}

func createChainHeight(d sqlx.DBExecutor, projectName string, r *CreateChainHeightReq, idg confid.SFIDGenerator) (*models.ChainHeight, error) {
	if err := checkChainID(d, r.ChainID); err != nil {
		return nil, err
	}

	n := *r
	n.EventType = types.EVENTTYPEDEFAULT // TODO support event type
	m := &models.ChainHeight{
		RelChainHeight: models.RelChainHeight{ChainHeightID: idg.MustGenSFID()},
		ChainHeightData: models.ChainHeightData{
			ProjectName:     projectName,
			ChainHeightInfo: n,
		},
	}
	if err := m.Create(d); err != nil {
		return nil, status.CheckDatabaseError(err, "CreateChainHeightMonitor")
	}
	return m, nil
}

func checkChainID(d sqlx.DBExecutor, id uint64) error {
	b := &models.Blockchain{RelBlockchain: models.RelBlockchain{ChainID: id}}
	if err := b.FetchByChainID(d); err != nil {
		return status.CheckDatabaseError(err, "GetBlockchainByChainID")
	}
	return nil
}
