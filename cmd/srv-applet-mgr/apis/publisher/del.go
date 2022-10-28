package publisher

import (
	"context"

	"github.com/iotexproject/Bumblebee/kit/httptransport/httpx"

	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/iotexproject/w3bstream/pkg/modules/publisher"
)

type RemovePublisher struct {
	httpx.MethodDelete
	publisher.RemovePublisherReq
}

func (r *RemovePublisher) Path() string { return "/:projectName" }

func (r *RemovePublisher) Output(ctx context.Context) (interface{}, error) {
	a := middleware.CurrentAccountFromContext(ctx)
	if _, err := a.ValidateProjectPermByPrjName(ctx, r.ProjectName); err != nil {
		return nil, err
	}

	return nil, publisher.RemovePublisher(ctx, &r.RemovePublisherReq)
}
