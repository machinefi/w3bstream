package blockchain

import (
	"context"

	confid "github.com/iotexproject/Bumblebee/conf/id"

	"github.com/iotexproject/w3bstream/pkg/errors/status"

	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type CreateContractlogReq = models.ContractlogData

func CreateContractlog(ctx context.Context, r *CreateContractlogReq) (*models.Contractlog, error) {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "CreateContractLog")
	defer l.End()

	n := *r
	n.BlockCurrent = n.BlockStart
	m := &models.Contractlog{
		RelContractlog:  models.RelContractlog{ContractlogID: idg.MustGenSFID()},
		ContractlogData: n,
	}
	if err := m.Create(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err)
	}

	return m, nil
}

type CreateChaintxReq = models.ChaintxData

func CreateChaintx(ctx context.Context, r *CreateChaintxReq) (*models.Chaintx, error) {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "CreateChainTx")
	defer l.End()

	m := &models.Chaintx{
		RelChaintx:  models.RelChaintx{ChaintxID: idg.MustGenSFID()},
		ChaintxData: *r,
	}
	if err := m.Create(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err)
	}
	return m, nil
}

type CreateChainHeightReq = models.ChainHeightData

func CreateChainHeight(ctx context.Context, r *CreateChainHeightReq) (*models.ChainHeight, error) {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "CreateChainHeight")
	defer l.End()

	m := &models.ChainHeight{
		RelChainHeight:  models.RelChainHeight{ChainHeightID: idg.MustGenSFID()},
		ChainHeightData: *r,
	}
	if err := m.Create(d); err != nil {
		return nil, status.CheckDatabaseError(err)
	}
	return m, nil
}
