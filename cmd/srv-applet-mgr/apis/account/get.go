package account

import (
	"context"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
)

// Deprecated use operator.ListOperator
type GetOperatorAddr struct {
	httpx.MethodGet `summary:"Get account operator by name"`

	AccountOperatorName string `in:"query" name:"accountOperatorName,omitempty"` // account operator name
}

func (r *GetOperatorAddr) Path() string { return "/operatoraddr" }

func (r *GetOperatorAddr) Output(ctx context.Context) (interface{}, error) {
	if r.AccountOperatorName == "" {
		r.AccountOperatorName = operator.DefaultOperatorName
	}

	ca := middleware.MustCurrentAccountFromContext(ctx)
	op, err := operator.GetByAccountAndName(ctx, ca.AccountID, r.AccountOperatorName)
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
