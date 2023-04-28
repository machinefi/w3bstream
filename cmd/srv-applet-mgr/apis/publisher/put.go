package publisher

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	"github.com/machinefi/w3bstream/pkg/types"
)

// Update Publisher by Publisher ID
type UpdatePublisher struct {
	httpx.MethodPut
	PublisherID         types.SFID `in:"path" name:"publisherID"`
	publisher.UpdateReq `in:"body"`
}

func (r *UpdatePublisher) Path() string { return "/:publisherID" }

func (r *UpdatePublisher) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithPublisherBySFID(ctx, r.PublisherID)
	if err != nil {
		return nil, err
	}
	r.PublisherID = types.MustPublisherFromContext(ctx).PublisherID

	return nil, publisher.Update(ctx, &r.UpdateReq)
}
