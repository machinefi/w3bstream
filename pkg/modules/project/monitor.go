package project

import (
	"context"

	confid "github.com/iotexproject/Bumblebee/conf/id"

	"github.com/iotexproject/w3bstream/pkg/enums"
	"github.com/iotexproject/w3bstream/pkg/errors/status"
	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type CreateMonitorReq struct {
	CreateContractlogReq `name:"contractlog"  json:"contractlog"`
}

type CreateContractlogReq = models.ContractlogInfo

func CreateMonitor(ctx context.Context, projectName string, r *CreateMonitorReq) (*models.Contractlog, error) {
	d := types.MustDBExecutorFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	b := &models.Blockchain{RelBlockchain: models.RelBlockchain{ChainID: r.ChainID}}
	if err := b.FetchByChainID(d); err != nil {
		return nil, status.CheckDatabaseError(err, "GetBlockchainByChainID")
	}

	n := r.CreateContractlogReq
	n.BlockCurrent = n.BlockStart
	n.EventType = enums.EVENT_TYPE__ANY // TODO support event type
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
