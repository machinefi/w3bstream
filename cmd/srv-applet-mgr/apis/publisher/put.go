package publisher

import (
	"context"

	"github.com/iotexproject/Bumblebee/kit/httptransport/httpx"

	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/iotexproject/w3bstream/pkg/modules/publisher"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type UpdatePublisher struct {
	httpx.MethodPut
	ProjectName                  string     `in:"path" name:"projectName"`
	PublisherID                  types.SFID `in:"path" name:"publisherID"`
	publisher.CreatePublisherReq `in:"body"`
}

func (r *UpdatePublisher) Path() string {
	return "/:projectName/:publisherID"
}

func (r *UpdatePublisher) Output(ctx context.Context) (interface{}, error) {
	a := middleware.CurrentAccountFromContext(ctx)
	if _, err := a.ValidateProjectPermByPrjName(ctx, r.ProjectName); err != nil {
		return nil, err
	}

	return nil, publisher.UpdatePublisher(ctx, r.PublisherID, &r.CreatePublisherReq)
}
