package publisher

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
)

type ListPublisher struct {
	httpx.MethodGet
	ProjectID types.SFID `name:"projectID" in:"path"`
	publisher.ListPublisherReq
}

func (r *ListPublisher) Path() string { return "/:projectID" }

func (r *ListPublisher) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)

	if _, err := ca.ValidateProjectPerm(ctx, r.ProjectID); err != nil {
		return nil, err
	}

	r.SetCurrentProject(r.ProjectID)
	return publisher.ListPublisher(ctx, &r.ListPublisherReq)
}
