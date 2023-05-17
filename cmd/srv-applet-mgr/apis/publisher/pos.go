package publisher

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	"github.com/machinefi/w3bstream/pkg/types"
)

// Create Publisher
type CreatePublisher struct {
	httpx.MethodPost
	publisher.CreateReq `in:"body"`
}

func (r *CreatePublisher) Output(ctx context.Context) (interface{}, error) {
	acc := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := acc.WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}

	prj := types.MustProjectFromContext(ctx)
	return publisher.Create(ctx, &acc.Account, prj, &r.CreateReq)
}
