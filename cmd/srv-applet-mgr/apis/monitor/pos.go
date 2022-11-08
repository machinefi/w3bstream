package monitor

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
)

type CreateContractLog struct {
	httpx.MethodPost
	ProjectName                     string `in:"path" name:"projectName"`
	blockchain.CreateContractLogReq `in:"body"`
}

func (r *CreateContractLog) Path() string { return "/contract_log/:projectName" }

func (r *CreateContractLog) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	_, err := ca.ValidateProjectPermByPrjName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	return blockchain.CreateContractLog(ctx, r.ProjectName, &r.CreateContractLogReq)
}

type CreateChainTx struct {
	httpx.MethodPost
	ProjectName                 string `in:"path" name:"projectName"`
	blockchain.CreateChainTxReq `in:"body"`
}

func (r *CreateChainTx) Path() string { return "/chain_tx/:projectName" }

func (r *CreateChainTx) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	_, err := ca.ValidateProjectPermByPrjName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	return blockchain.CreateChainTx(ctx, r.ProjectName, &r.CreateChainTxReq)
}

type CreateChainHeight struct {
	httpx.MethodPost
	ProjectName                     string `in:"path" name:"projectName"`
	blockchain.CreateChainHeightReq `in:"body"`
}

func (r *CreateChainHeight) Path() string { return "/chain_height/:projectName" }

func (r *CreateChainHeight) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	_, err := ca.ValidateProjectPermByPrjName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	return blockchain.CreateChainHeight(ctx, r.ProjectName, &r.CreateChainHeightReq)
}
