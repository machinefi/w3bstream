package projectoperator

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/projectoperator"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Create struct {
	httpx.MethodPost
	ProjectID  types.SFID `in:"path" name:"projectID"`
	OperatorID types.SFID `in:"path" name:"operatorID"`
}

func (r *Create) Path() string { return "/:projectID/:operatorID" }

func (r *Create) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithOperatorBySFID(ctx, r.OperatorID)
	if err != nil {
		return nil, err
	}
	ctx, err = middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextBySFID(ctx, r.ProjectID)
	if err != nil {
		return nil, err
	}

	return projectoperator.Create(ctx, r.ProjectID, r.OperatorID)
}
