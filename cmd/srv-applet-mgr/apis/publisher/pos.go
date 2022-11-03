package publisher

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
)

type CreatePublisher struct {
	httpx.MethodPost
	ProjectID                    types.SFID `in:"path" name:"projectID"`
	publisher.CreatePublisherReq `in:"body"`
}

func (r *CreatePublisher) Path() string {
	return "/:projectID"
}

func (r *CreatePublisher) Output(ctx context.Context) (interface{}, error) {
	a := middleware.CurrentAccountFromContext(ctx)
	if _, err := a.ValidateProjectPerm(ctx, r.ProjectID); err != nil {
		return nil, err
	}

	return publisher.CreatePublisher(ctx, r.ProjectID, &r.CreatePublisherReq)
}
