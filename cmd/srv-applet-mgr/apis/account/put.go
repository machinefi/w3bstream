package account

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/account"
)

type UpdatePasswordByAccountID struct {
	httpx.MethodPut `summary:"Update account password"`

	account.UpdatePasswordReq `in:"body"`
}

func (r *UpdatePasswordByAccountID) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	return nil, account.UpdateAccountPassword(ctx, ca.AccountID, &r.UpdatePasswordReq)
}
