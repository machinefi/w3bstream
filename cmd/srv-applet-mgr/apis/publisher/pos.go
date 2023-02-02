package publisher

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreatePublisher struct {
	httpx.MethodPost
	ProjectName                  string `in:"path" name:"projectName"`
	publisher.CreatePublisherReq `in:"body"`
}

func (r *CreatePublisher) Path() string {
	return "/:projectName"
}

func (r *CreatePublisher) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)

	return publisher.CreatePublisher(ctx, prj.ProjectID, &r.CreatePublisherReq)
}
