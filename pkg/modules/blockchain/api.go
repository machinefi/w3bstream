package blockchain

import (
	"context"
	"errors"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

const chainUniqFlag = 0

type CreateMonitorReq struct {
	Contractlog *CreateContractlogReq `json:"contractLog,omitempty"`
	Chaintx     *CreateChaintxReq     `json:"chainTx,omitempty"`
	ChainHeight *CreateChainHeightReq `json:"chainHeight,omitempty"`
}

type (
	CreateContractlogReq = models.ContractLogInfo
	CreateChaintxReq     = models.ChainTxInfo
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

func createContractLog(d sqlx.DBExecutor, projectName string, r *CreateContractlogReq, idg confid.SFIDGenerator) (*models.ContractLog, error) {
	if err := checkChainID(d, r.ChainID); err != nil {
		return nil, err
	}

	n := *r
	n.BlockCurrent = n.BlockStart
	n.EventType = getEventType(n.EventType)
	m := &models.ContractLog{
		RelContractlog: models.RelContractlog{ContractLogID: idg.MustGenSFID()},
		ContractLogData: models.ContractLogData{
			ProjectName:     projectName,
			Uniq:            chainUniqFlag,
			ContractLogInfo: n,
		},
	}
	if err := m.Create(d); err != nil {
		return nil, status.CheckDatabaseError(err, "CreateContractlogMonitor")
	}
	return m, nil
}

func createChainTx(d sqlx.DBExecutor, projectName string, r *CreateChaintxReq, idg confid.SFIDGenerator) (*models.ChainTx, error) {
	if err := checkChainID(d, r.ChainID); err != nil {
		return nil, err
	}

	n := *r
	n.EventType = getEventType(n.EventType)
	m := &models.ChainTx{
		RelChainTx: models.RelChainTx{ChainTxID: idg.MustGenSFID()},
		ChainTxData: models.ChainTxData{
			ProjectName: projectName,
			Uniq:        chainUniqFlag,
			Finished:    datatypes.FALSE,
			ChainTxInfo: n,
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
	n.EventType = getEventType(n.EventType)
	m := &models.ChainHeight{
		RelChainHeight: models.RelChainHeight{ChainHeightID: idg.MustGenSFID()},
		ChainHeightData: models.ChainHeightData{
			ProjectName:     projectName,
			Uniq:            chainUniqFlag,
			Finished:        datatypes.FALSE,
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

type RemoveMonitorReq struct {
	ContractlogID types.SFID `json:"contractlogID,omitempty"`
	ChaintxID     types.SFID `json:"chaintxID,omitempty"`
	ChainHeightID types.SFID `json:"chainHeightID,omitempty"`
}

func RemoveMonitor(ctx context.Context, projectName string, r *RemoveMonitorReq) error {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "RemoveMonitor")
	defer l.End()

	l = l.WithValues("project", projectName)

	switch {
	case r.ContractlogID != 0:
		m := &models.ContractLog{RelContractlog: models.RelContractlog{ContractLogID: r.ContractlogID}}
		if err := m.FetchByContractLogID(d); err != nil {
			return status.CheckDatabaseError(err, "FetchByContractlogID")
		}
		if err := checkProjectName(m.ProjectName, projectName, l); err != nil {
			return err
		}
		if err := m.DeleteByContractLogID(d); err != nil {
			return status.CheckDatabaseError(err, "DeleteByContractlogID")
		}

	case r.ChaintxID != 0:
		m := &models.ChainTx{RelChainTx: models.RelChainTx{ChainTxID: r.ChaintxID}}
		if err := m.FetchByChainTxID(d); err != nil {
			return status.CheckDatabaseError(err, "FetchByChaintxID")
		}
		if err := checkProjectName(m.ProjectName, projectName, l); err != nil {
			return err
		}
		if err := m.DeleteByChainTxID(d); err != nil {
			return status.CheckDatabaseError(err, "DeleteByChaintxID")
		}

	case r.ChainHeightID != 0:
		m := &models.ChainHeight{RelChainHeight: models.RelChainHeight{ChainHeightID: r.ChainHeightID}}
		if err := m.FetchByChainHeightID(d); err != nil {
			return status.CheckDatabaseError(err, "FetchByChainHeightID")
		}
		if err := checkProjectName(m.ProjectName, projectName, l); err != nil {
			return err
		}
		if err := m.DeleteByChainHeightID(d); err != nil {
			return status.CheckDatabaseError(err, "DeleteByChainHeightID")
		}

	default:
		return status.BadRequest
	}

	return nil
}

func checkProjectName(want, curr string, l log.Logger) error {
	if want != curr {
		l.Error(errors.New("monitor project mismatch"))
		return status.BadRequest.StatusErr().WithDesc("monitor project mismatch")
	}
	return nil
}

func getEventType(eventType string) string {
	if eventType == "" {
		return enums.MONITOR_EVENTTYPEDEFAULT
	}
	return eventType
}
