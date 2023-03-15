package login

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/account"
)

type GetNonceByEthAddress struct {
	httpx.MethodGet
	Address types.EthAddress `in:"path" name:"address" validate:"@ethAddress"`
}

func (r *GetNonceByEthAddress) Path() string {
	return "/nonce/:ethAddress"
}

func (r *GetNonceByEthAddress) Output(ctx context.Context) (interface{}, error) {
	return account.GetNonce(ctx, r.Address)
}
