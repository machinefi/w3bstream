package account

import (
	"context"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
)

type GetOperatorAddr struct {
	httpx.MethodGet
}

func (r *GetOperatorAddr) Path() string { return "/operatoraddr" }

func (r *GetOperatorAddr) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	prvkey, err := crypto.HexToECDSA(ca.OperatorPrivateKey)
	if err != nil {
		return nil, err
	}
	pubkey := crypto.PubkeyToAddress(prvkey.PublicKey)
	return pubkey.Hex(), nil
}
