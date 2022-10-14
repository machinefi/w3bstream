package blockchain

import (
	"context"

	confid "github.com/iotexproject/Bumblebee/conf/id"

	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type CreateContractlogReq = models.ContractlogInfo

func CreateContractlog(ctx context.Context, r *CreateContractlogReq) (*models.Contractlog, error) {
	d := types.MustDBExecutorFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	n := *r
	n.BlockCurrent = n.BlockStart
	m := &models.Contractlog{
		RelContractlog:  models.RelContractlog{ContractlogID: idg.MustGenSFID()},
		ContractlogInfo: n,
	}

	if err := m.Create(d); err != nil {
		return nil, err
	}

	return m, nil
}
