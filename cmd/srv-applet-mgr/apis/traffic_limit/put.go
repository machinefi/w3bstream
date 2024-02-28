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
	_, ok := middleware.MustCurrentAccountFromContext(ctx).CheckRole(enums.ACCOUNT_ROLE__ADMIN)
	if !ok {
		return nil, status.NoAdminPermission
	}

	return trafficlimit.Update(ctx, r.TrafficLimitID, &r.UpdateReq)
}
