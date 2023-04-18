package deploy

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	"github.com/machinefi/w3bstream/pkg/types"
)

type RemoveInstance struct {
	httpx.MethodDelete
	InstanceID types.SFID `in:"path" name:"instanceID"`
}

func (r *RemoveInstance) Path() string { return "/:instanceID" }

func (r *RemoveInstance) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.CurrentAccountFromContext(ctx).
		WithInstanceContextBySFID(ctx, r.InstanceID)
	if err != nil {
		return nil, err
	}
	return nil, deploy.Remove(ctx, false)
}
