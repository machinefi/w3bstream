package projectoperator

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
	"github.com/machinefi/w3bstream/pkg/modules/projectoperator"
	"github.com/machinefi/w3bstream/pkg/types"
)

type GetProjectOperator struct {
	httpx.MethodGet
	ProjectID types.SFID `in:"path" name:"projectID"`
}

func (r *GetProjectOperator) Path() string { return "/data/:projectID" }

func (r *GetProjectOperator) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextBySFID(ctx, r.ProjectID)
	if err != nil {
		return nil, err
	}

	po, err := projectoperator.GetByProject(ctx, r.ProjectID)
	if err != nil {
		if err == status.ProjectOperatorNotFound {
			return operator.GetDetailByAccountAndName(ctx, ca.AccountID, operator.DefaultOperatorName)
		}
		return nil, err
	}
	return operator.GetDetailBySFID(ctx, po.OperatorID)
}
