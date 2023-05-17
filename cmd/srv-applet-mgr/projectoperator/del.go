package projectoperator

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/projectoperator"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Remove struct {
	httpx.MethodDelete
	ProjectID types.SFID `in:"path" name:"projectID"`
}

func (r *Remove) Path() string { return "/:projectID" }

func (r *Remove) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextBySFID(ctx, r.ProjectID)
	if err != nil {
		return nil, err
	}
	return nil, projectoperator.RemoveByProject(ctx, r.ProjectID)
}
