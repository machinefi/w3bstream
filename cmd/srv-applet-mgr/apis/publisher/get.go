package publisher

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	"github.com/machinefi/w3bstream/pkg/types"
)

type ListPublisher struct {
	httpx.MethodGet
	ProjectName string `in:"path" name:"projectName"`
	publisher.ListPublisherReq
}

func (r *ListPublisher) Path() string { return "/:projectName" }

func (r *ListPublisher) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)

	r.SetCurrentProject(prj.ProjectID)
	return publisher.ListPublisher(ctx, &r.ListPublisherReq)
}
