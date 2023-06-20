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

type ListTrafficLimit struct {
	httpx.MethodGet
	trafficlimit.ListReq
}

func (r *ListTrafficLimit) Path() string { return "/datalist" }

func (r *ListTrafficLimit) Output(ctx context.Context) (interface{}, error) {
	ca, ok := middleware.MustCurrentAccountFromContext(ctx).CheckRole(enums.ACCOUNT_ROLE__ADMIN)
	if !ok {
		return nil, status.NoAdminPermission
	}

	ctx, err := ca.WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}

	r.ProjectID = types.MustProjectFromContext(ctx).ProjectID
	return trafficlimit.List(ctx, &r.ListReq)
}

type GetTrafficLimit struct {
	httpx.MethodGet
	TrafficLimitID types.SFID `in:"path" name:"trafficLimitID"`
}

func (r *GetTrafficLimit) Path() string { return "/data/:trafficLimitID" }

func (r *GetTrafficLimit) Output(ctx context.Context) (interface{}, error) {
	ca, ok := middleware.MustCurrentAccountFromContext(ctx).CheckRole(enums.ACCOUNT_ROLE__ADMIN)
	if !ok {
		return nil, status.NoAdminPermission
	}

	ctx, err := ca.WithTrafficLimitContextBySFID(ctx, r.TrafficLimitID)
	if err != nil {
		return nil, err
	}

	return types.MustTrafficLimitFromContext(ctx), nil
}
