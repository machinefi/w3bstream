package publisher

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	"github.com/machinefi/w3bstream/pkg/types"
)

// Remove Publisher by Publisher ID
type RemovePublisher struct {
	httpx.MethodDelete
	PublisherID types.SFID `in:"path" name:"publisherID"`
}

func (r *RemovePublisher) Path() string { return "/data/:publisherID" }

func (r *RemovePublisher) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithPublisherBySFID(ctx, r.PublisherID)
	if err != nil {
		return nil, err
	}
	return nil, publisher.RemoveBySFID(ctx, r.PublisherID)
}

// Remove Publisher by Given Conditions
type BatchRemovePublisher struct {
	httpx.MethodDelete
	publisher.CondArgs
}

func (r *BatchRemovePublisher) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}
	r.ProjectIDs = []types.SFID{types.MustProjectFromContext(ctx).ProjectID}

	return nil, publisher.Remove(ctx, &r.CondArgs)
}
