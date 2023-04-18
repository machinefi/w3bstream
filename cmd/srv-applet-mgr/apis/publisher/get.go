package publisher

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	"github.com/machinefi/w3bstream/pkg/types"
)

type ListPublisher struct {
	httpx.MethodGet
	publisher.ListPublisherReq
}

func (r *ListPublisher) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)

	r.SetCurrentProject(prj.ProjectID)
	return publisher.ListPublisher(ctx, &r.ListPublisherReq)
}
