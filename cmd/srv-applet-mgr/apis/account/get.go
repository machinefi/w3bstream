package account

import (
	"context"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
)

type GetOperatorAddr struct {
	httpx.MethodGet
}

func (r *GetOperatorAddr) Path() string { return "/operatoraddr" }

func (r *GetOperatorAddr) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	op, err := operator.GetByAccountAndName(ctx, ca.AccountID, operator.DefaultOperatorName)
	if err != nil {
		return nil, err
	}
	prvkey, err := crypto.HexToECDSA(op.PrivateKey)
	if err != nil {
		return nil, err
	}
	pubkey := crypto.PubkeyToAddress(prvkey.PublicKey)
	return pubkey.Hex(), nil
}
