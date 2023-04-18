package publisher

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	"github.com/machinefi/w3bstream/pkg/types"
)

type RemovePublisher struct {
	httpx.MethodDelete
	PublisherID types.SFID `in:"path" name:"publisherID"`
}

func (r *RemovePublisher) Path() string { return "/data/:publisherID" }

func (r *RemovePublisher) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithPublisherBySFID(ctx, r.PublisherID)
	if err != nil {
		return nil, err
	}
	return nil, publisher.Remove(ctx)
}

type BatchRemoveByPublisherIDs struct {
	httpx.MethodDelete
	PublisherIDs []types.SFID `in:"query" name:"publisherID"`
}

func (r *BatchRemoveByPublisherIDs) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, middleware.ProjectProviderFromContext(ctx))
	if err != nil {
		return nil, err
	}
	return nil, publisher.BatchRemoveBySFIDs(ctx, r.PublisherIDs)
}
