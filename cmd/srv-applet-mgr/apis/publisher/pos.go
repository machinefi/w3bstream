package publisher

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
)

type CreatePublisher struct {
	httpx.MethodPost
	publisher.CreateReq `in:"body"`
}

func (r *CreatePublisher) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, middleware.ProjectProviderFromContext(ctx))
	if err != nil {
		return nil, err
	}

	return publisher.Create(ctx, &r.CreateReq)
}
