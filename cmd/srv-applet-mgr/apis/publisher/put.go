package publisher

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	"github.com/machinefi/w3bstream/pkg/types"
)

type UpdatePublisher struct {
	httpx.MethodPut
	PublisherID         types.SFID `in:"path" name:"publisherID"`
	publisher.CreateReq `in:"body"`
}

func (r *UpdatePublisher) Path() string { return "/:publisherID" }

func (r *UpdatePublisher) Output(ctx context.Context) (interface{}, error) {
	a := middleware.CurrentAccountFromContext(ctx)
	ctx, err := a.WithPublisherBySFID(ctx, r.PublisherID)
	if err != nil {
		return nil, err
	}
	return nil, publisher.Update(ctx, &r.CreateReq)
}
