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
	PublisherID                  types.SFID `in:"path" name:"publisherID"`
	publisher.CreatePublisherReq `in:"body"`
}

func (r *UpdatePublisher) Path() string { return "/:publisherID" }

func (r *UpdatePublisher) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).WithPublisherBySFID(ctx, r.PublisherID)
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)

	return nil, publisher.UpdatePublisher(ctx, prj, r.PublisherID, &r.CreatePublisherReq)
}
