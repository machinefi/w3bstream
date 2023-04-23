package publisher

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	"github.com/machinefi/w3bstream/pkg/types"
)

// Get Publisher by Publisher ID
type GetPublisher struct {
	httpx.MethodGet
	PublisherID types.SFID `in:"path" name:"publisherID"`
}

func (r *GetPublisher) Path() string { return "/data/:publisherID" }

func (r *GetPublisher) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithPublisherBySFID(ctx, r.PublisherID)
	if err != nil {
		return nil, err
	}
	return types.MustPublisherFromContext(ctx), nil
}

// List Publishers by Conditions
type ListPublisher struct {
	httpx.MethodGet
	publisher.ListReq
}

func (r *ListPublisher) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}
	r.PublisherIDs = []types.SFID{types.MustProjectFromContext(ctx).ProjectID}
	return publisher.List(ctx, &r.ListReq)
}
