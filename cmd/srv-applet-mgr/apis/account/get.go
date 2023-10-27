package account

import (
	"context"
	"crypto/ed25519"

	solcommon "github.com/blocto/solana-go-sdk/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/enums"
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
	if op.Type == enums.OPERATOR_KEY__ECDSA {
		prvkey, err := crypto.HexToECDSA(op.PrivateKey)
		if err != nil {
			return nil, err
		}
		pubkey := crypto.PubkeyToAddress(prvkey.PublicKey)
		return pubkey.Hex(), nil
	} else {
		b := common.FromHex(op.PrivateKey)
		prik := ed25519.PrivateKey(b)
		pubk := solcommon.PublicKeyFromBytes(prik.Public().(ed25519.PublicKey))
		return pubk.ToBase58(), nil
	}
}
