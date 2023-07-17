package account_access_key

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
)

// DeleteAccountAccessKeyByName delete access key by name under current account
type DeleteAccountAccessKeyByName struct {
	httpx.MethodDelete
	Name string `in:"path" name:"name"`
}

func (r *DeleteAccountAccessKeyByName) Path() string { return "/:name" }

func (r *DeleteAccountAccessKeyByName) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	return nil, access_key.DeleteByName(ca.WithAccount(ctx), r.Name)
}
