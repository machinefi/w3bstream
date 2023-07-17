package account_access_key

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
)

// CreateAccountAccessKey create account access key
type CreateAccountAccessKey struct {
	httpx.MethodPost
	access_key.CreateAccountAccessKeyReq `in:"body"`
}

func (r *CreateAccountAccessKey) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	return access_key.Create(ca.WithAccount(ctx), &access_key.CreateReq{
		IdentityID:    ca.AccountID,
		IdentityType:  enums.ACCESS_KEY_IDENTITY_TYPE__ACCOUNT,
		CreateReqBase: r.CreateAccountAccessKeyReq,
	})
}
