package deploy

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	"github.com/machinefi/w3bstream/pkg/types"
)

// RemoveInstance remove instance by instance id
type RemoveInstance struct {
	httpx.MethodDelete
	InstanceID types.SFID `in:"path" name:"instanceID"`
}

func (r *RemoveInstance) Path() string { return "/data/:instanceID" }

func (r *RemoveInstance) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithInstanceContextBySFID(ctx, r.InstanceID)
	if err != nil {
		return nil, err
	}

	return nil, deploy.RemoveBySFID(ctx, r.InstanceID)
}

// BatchRemoveInstance remove instances by condition
type BatchRemoveInstance struct {
	httpx.MethodDelete
	deploy.CondArgs
}

func (r *BatchRemoveInstance) Path() string { return "" }

func (r *BatchRemoveInstance) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}
	r.CondArgs.ProjectID = types.MustProjectFromContext(ctx).ProjectID
	return nil, deploy.Remove(ctx, &r.CondArgs)
}
