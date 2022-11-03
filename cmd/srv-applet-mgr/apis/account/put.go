package account

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/account"
)

type UpdatePasswordByAccountID struct {
	httpx.MethodPut
	AccountID                 types.SFID `in:"path" name:"accountID"`
	account.UpdatePasswordReq `in:"body"`
}

func (r *UpdatePasswordByAccountID) Path() string { return "/:accountID" }

func (r *UpdatePasswordByAccountID) Output(ctx context.Context) (interface{}, error) {
	return nil, account.UpdateAccountPassword(ctx, r.AccountID, r.Password)
}
