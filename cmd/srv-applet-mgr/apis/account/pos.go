package account

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/account"
)

type CreateAccountByUsernameAndPassword struct {
	httpx.MethodPost
	account.CreateAccountByUsernameReq `in:"body"`
}

func (r *CreateAccountByUsernameAndPassword) Path() string { return "/admin" }

func (r *CreateAccountByUsernameAndPassword) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	if ca.Role != enums.ACCOUNT_ROLE__ADMIN {
		return nil, status.NoAdminPermission
	}
	return account.CreateAccountByUsername(ctx, &r.CreateAccountByUsernameReq)
}
