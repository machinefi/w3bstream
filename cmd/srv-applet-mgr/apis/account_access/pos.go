package account_access

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/account_access"
)

type CreateAccountAccessKey struct {
	httpx.MethodPost
	account_access.CreateReq `in:"body"`
}

func (r *CreateAccountAccessKey) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	return account_access.Create(ca.WithAccount(ctx), &r.CreateReq)
}
