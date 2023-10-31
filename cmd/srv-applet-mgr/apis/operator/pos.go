package operator

import (
	"context"
	"crypto/ed25519"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
)

type CreateOperator struct {
	httpx.MethodPost
	operator.CreateReq `in:"body"`
}

func (r *CreateOperator) Output(ctx context.Context) (interface{}, error) {
	ctx = middleware.MustCurrentAccountFromContext(ctx).WithAccount(ctx)

	if r.Type == enums.OPERATOR_KEY__ECDSA {
		if _, err := crypto.HexToECDSA(r.PrivateKey); err != nil {
			return nil, status.InvalidPrivateKey.StatusErr().WithDesc(err.Error())
		}
	} else {
		b := common.FromHex(r.PrivateKey)
		if len(b) != ed25519.PrivateKeySize {
			return nil, status.InvalidPrivateKey.StatusErr().WithDesc(fmt.Sprintf("invalid key size, expect %d, get %d", ed25519.PrivateKeySize, len(b)))
		}
	}

	return operator.Create(ctx, &r.CreateReq)
}
