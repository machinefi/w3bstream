package resource

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
)

type ListResources struct {
	httpx.MethodGet
	resource.ListReq
}

func (r *ListResources) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)

	r.AccountID = ca.AccountID
	return resource.List(ctx, &r.ListReq)
}
