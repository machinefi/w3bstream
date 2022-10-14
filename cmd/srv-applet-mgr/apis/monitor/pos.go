package monitor

import (
	"context"

	"github.com/iotexproject/Bumblebee/kit/httptransport/httpx"

	"github.com/iotexproject/w3bstream/pkg/modules/blockchain"
)

type CreateContractlog struct {
	httpx.MethodPost
	blockchain.CreateContractlogReq `in:"body"`
}

func (r *CreateContractlog) Output(ctx context.Context) (interface{}, error) {
	return blockchain.CreateContractlog(ctx, &r.CreateContractlogReq)
}
