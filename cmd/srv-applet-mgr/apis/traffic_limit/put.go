package traffic_limit

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/trafficlimit"
	"github.com/machinefi/w3bstream/pkg/types"
)

type UpdateTrafficLimit struct {
	httpx.MethodPut
	TrafficLimitID         types.SFID `in:"path" name:"trafficLimitID"`
	trafficlimit.UpdateReq `in:"body"`
}

func (r *UpdateTrafficLimit) Path() string {
	return "/:trafficLimitID"
}

func (r *UpdateTrafficLimit) Output(ctx context.Context) (interface{}, error) {
	ca, ok := middleware.MustCurrentAccountFromContext(ctx).CheckRole(enums.ACCOUNT_ROLE__ADMIN)
	if !ok {
		return nil, status.NoAdminPermission
	}

	ctx, err := ca.WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}

	ctx, err = ca.WithTrafficLimitContextBySFIDAndProjectName(ctx, r.TrafficLimitID)
	if err != nil {
		return nil, err
	}

	r.UpdateReq.TrafficLimitID = r.TrafficLimitID
	return trafficlimit.Update(ctx, &r.UpdateReq)
}
