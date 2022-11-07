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
	ProjectID                       types.SFID `in:"path" name:"projectID"`
	blockchain.CreateContractLogReq `in:"body"`
}

func (r *CreateContractLog) Path() string { return "/contract_log/:projectID" }

func (r *CreateContractLog) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	p, err := ca.ValidateProjectPerm(ctx, r.ProjectID)
	if err != nil {
		return nil, err
	}
	return blockchain.CreateContractLog(ctx, p.Name, &r.CreateContractLogReq)
}

type CreateChainTx struct {
	httpx.MethodPost
	ProjectID                   types.SFID `in:"path" name:"projectID"`
	blockchain.CreateChainTxReq `in:"body"`
}

func (r *CreateChainTx) Path() string { return "/chain_tx/:projectID" }

func (r *CreateChainTx) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	p, err := ca.ValidateProjectPerm(ctx, r.ProjectID)
	if err != nil {
		return nil, err
	}
	return blockchain.CreateChainTx(ctx, p.Name, &r.CreateChainTxReq)
}

type CreateChainHeight struct {
	httpx.MethodPost
	ProjectID                       types.SFID `in:"path" name:"projectID"`
	blockchain.CreateChainHeightReq `in:"body"`
}

func (r *CreateChainHeight) Path() string { return "/chain_height/:projectID" }

func (r *CreateChainHeight) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	p, err := ca.ValidateProjectPerm(ctx, r.ProjectID)
	if err != nil {
		return nil, err
	}
	return blockchain.CreateChainHeight(ctx, p.Name, &r.CreateChainHeightReq)
}
