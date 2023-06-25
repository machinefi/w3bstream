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

// RemoveTrafficLimit Remove TrafficLimit by TrafficLimit ID
type RemoveTrafficLimit struct {
	httpx.MethodDelete
	TrafficLimitID types.SFID `in:"path" name:"trafficLimitID"`
}

func (r *RemoveTrafficLimit) Path() string { return "/data/:trafficLimitID" }

func (r *RemoveTrafficLimit) Output(ctx context.Context) (interface{}, error) {
	acc, ok := middleware.MustCurrentAccountFromContext(ctx).CheckRole(enums.ACCOUNT_ROLE__ADMIN)
	if !ok {
		return nil, status.NoAdminPermission
	}

	ctx, err := acc.WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}
	return nil, trafficlimit.RemoveBySFID(ctx, r.TrafficLimitID)
}

// BatchRemoveTrafficLimit Remove TrafficLimit by Given Conditions
type BatchRemoveTrafficLimit struct {
	httpx.MethodDelete
	trafficlimit.CondArgs
}

func (r *BatchRemoveTrafficLimit) Output(ctx context.Context) (interface{}, error) {
	acc, ok := middleware.MustCurrentAccountFromContext(ctx).CheckRole(enums.ACCOUNT_ROLE__ADMIN)
	if !ok {
		return nil, status.NoAdminPermission
	}

	ctx, err := acc.WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}
	r.ProjectID = types.MustProjectFromContext(ctx).ProjectID
	return nil, trafficlimit.Remove(ctx, &r.CondArgs)
}
