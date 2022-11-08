package monitor

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
	"github.com/machinefi/w3bstream/pkg/types"
)

type RemoveContractLog struct {
	httpx.MethodDelete
	ProjectName   string     `in:"path" name:"projectName"`
	ContractLogID types.SFID `in:"path" name:"contractLogID"`
}

func (r *RemoveContractLog) Path() string { return "/contract_log/:projectName" }

func (r *RemoveContractLog) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	_, err := ca.ValidateProjectPermByPrjName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	return nil, blockchain.RemoveContractLog(ctx, r.ProjectName, r.ContractLogID)
}

type RemoveChainTx struct {
	httpx.MethodDelete
	ProjectName string     `in:"path" name:"projectName"`
	ChainTxID   types.SFID `in:"path" name:"chainTxID"`
}

func (r *RemoveChainTx) Path() string { return "/chain_tx/:projectName" }

func (r *RemoveChainTx) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	_, err := ca.ValidateProjectPermByPrjName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	return nil, blockchain.RemoveChainTx(ctx, r.ProjectName, r.ChainTxID)
}

type RemoveChainHeight struct {
	httpx.MethodDelete
	ProjectName   string     `in:"path" name:"projectName"`
	ChainHeightID types.SFID `in:"path" name:"chainHeightID"`
}

func (r *RemoveChainHeight) Path() string { return "/chain_height/:projectName" }

func (r *RemoveChainHeight) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	_, err := ca.ValidateProjectPermByPrjName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	return nil, blockchain.RemoveChainHeight(ctx, r.ProjectName, r.ChainHeightID)
}
