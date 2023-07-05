package account_access_key

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
)

// ListAccountAccessKey get api access key list under current account
type ListAccountAccessKey struct {
	httpx.MethodGet
	access_key.ListReq
}

func (r *ListAccountAccessKey) Path() string { return "/datalist" }

func (r *ListAccountAccessKey) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	r.AccountID = ca.AccountID
	return access_key.List(ctx, &r.ListReq)
}
