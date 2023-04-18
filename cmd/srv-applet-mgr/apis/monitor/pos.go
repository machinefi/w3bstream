package monitor

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
)

type CreateContractLog struct {
	httpx.MethodPost
	blockchain.CreateContractLogReq `in:"body"`
}

func (r *CreateContractLog) Path() string { return "/contract_log" }

func (r *CreateContractLog) Output(ctx context.Context) (interface{}, error) {
	prj := middleware.MustProjectName(ctx)
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, prj)
	if err != nil {
		return nil, err
	}
	return blockchain.CreateContractLog(ctx, prj, &r.CreateContractLogReq)
}

type CreateChainTx struct {
	httpx.MethodPost
	blockchain.CreateChainTxReq `in:"body"`
}

func (r *CreateChainTx) Path() string { return "/chain_tx" }

func (r *CreateChainTx) Output(ctx context.Context) (interface{}, error) {
	prj := middleware.MustProjectName(ctx)
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, prj)
	if err != nil {
		return nil, err
	}
	return blockchain.CreateChainTx(ctx, prj, &r.CreateChainTxReq)
}

type CreateChainHeight struct {
	httpx.MethodPost
	blockchain.CreateChainHeightReq `in:"body"`
}

func (r *CreateChainHeight) Path() string { return "/chain_height" }

func (r *CreateChainHeight) Output(ctx context.Context) (interface{}, error) {
	prj := middleware.MustProjectName(ctx)
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, prj)
	if err != nil {
		return nil, err
	}
	return blockchain.CreateChainHeight(ctx, prj, &r.CreateChainHeightReq)
}
