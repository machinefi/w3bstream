package account_access_key

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
)

type UpdateAccountAccessKeyByName struct {
	httpx.MethodPut `summary:"Update current account access key by name"`

	Name                 string `in:"path" name:"name"`
	access_key.UpdateReq `in:"body"`
}

func (r *UpdateAccountAccessKeyByName) Path() string { return "/:name" }

func (r *UpdateAccountAccessKeyByName) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	return access_key.UpdateByName(ca.WithAccount(ctx), r.Name, &r.UpdateReq)
}
