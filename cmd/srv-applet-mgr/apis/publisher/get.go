package publisher

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
)

type ListPublisher struct {
	httpx.MethodGet
	publisher.ListReq
}

func (r *ListPublisher) Path() string { return "/data_list" }

func (r *ListPublisher) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, middleware.ProjectProviderFromContext(ctx))
	if err != nil {
		return nil, err
	}
	return publisher.List(ctx, &r.ListReq)
}

type ListPublisherDetail struct {
	httpx.MethodGet
	publisher.ListReq
}

func (r *ListPublisherDetail) Path() string { return "/detail_list" }

func (r *ListPublisherDetail) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, middleware.ProjectProviderFromContext(ctx))
	if err != nil {
		return nil, err
	}
	return publisher.ListDetail(ctx, &r.ListReq)
}
