package account

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/account"
)

type CreateAccount struct {
	httpx.MethodPost
	account.CreateAccountByUsernameReq `in:"body"`
}

func (r *CreateAccount) Output(ctx context.Context) (interface{}, error) {
	return account.CreateAccountByUsername(ctx, &r.CreateAccountByUsernameReq)
}
