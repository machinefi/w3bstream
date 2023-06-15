package account_access

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/account_access"
)

type DeleteAccountAccessKeyByName struct {
	httpx.MethodDelete
	Name string `in:"path" name:"name"`
}

func (r *DeleteAccountAccessKeyByName) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	return nil, account_access.DeleteByName(ca.WithAccount(ctx), r.Name)
}
