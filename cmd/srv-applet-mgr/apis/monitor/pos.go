package monitor

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateContractLog struct {
	httpx.MethodPost
	blockchain.CreateContractLogReq `in:"body"`
}

func (r *CreateContractLog) Path() string { return "/contract_log" }

func (r *CreateContractLog) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}

	r.ProjectName = types.MustProjectFromContext(ctx).ProjectName.Name
	return blockchain.CreateContractLog(ctx, &r.CreateContractLogReq)
}

type CreateChainTx struct {
	httpx.MethodPost
	blockchain.CreateChainTxReq `in:"body"`
}

func (r *CreateChainTx) Path() string { return "/chain_tx" }

func (r *CreateChainTx) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}

	r.ProjectName = types.MustProjectFromContext(ctx).ProjectName.Name
	return blockchain.CreateChainTx(ctx, &r.CreateChainTxReq)
}

type CreateChainHeight struct {
	httpx.MethodPost
	blockchain.CreateChainHeightReq `in:"body"`
}

func (r *CreateChainHeight) Path() string { return "/chain_height" }

func (r *CreateChainHeight) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}

	r.ProjectName = types.MustProjectFromContext(ctx).ProjectName.Name
	return blockchain.CreateChainHeight(ctx, &r.CreateChainHeightReq)
}
