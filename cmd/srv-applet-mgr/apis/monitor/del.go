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
	ContractLogID types.SFID `in:"path" name:"contractLogID"`
}

func (r *RemoveContractLog) Path() string { return "/contract_log/:contractLogID" }

func (r *RemoveContractLog) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithContractLogBySFID(ctx, r.ContractLogID)
	if err != nil {
		return nil, err
	}

	return nil, blockchain.RemoveContractLogBySFID(ctx, r.ContractLogID)
}

type RemoveChainTx struct {
	httpx.MethodDelete
	ChainTxID types.SFID `in:"path" name:"chainTxID"`
}

func (r *RemoveChainTx) Path() string { return "/chain_tx/:chainTxID" }

func (r *RemoveChainTx) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithChainTxBySFID(ctx, r.ChainTxID)
	if err != nil {
		return nil, err
	}

	return nil, blockchain.RemoveChainTxBySFID(ctx, r.ChainTxID)
}

type RemoveChainHeight struct {
	httpx.MethodDelete
	ChainHeightID types.SFID `in:"path" name:"chainHeightID"`
}

func (r *RemoveChainHeight) Path() string { return "/chain_height/:chainHeightID" }

func (r *RemoveChainHeight) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithContractLogBySFID(ctx, r.ChainHeightID)
	if err != nil {
		return nil, err
	}

	return nil, blockchain.RemoveChainHeightBySFID(ctx, r.ChainHeightID)
}
