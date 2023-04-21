package deploy

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
)

type ControlInstance struct {
	httpx.MethodPut
	InstanceID types.SFID      `in:"path" name:"instanceID"`
	Cmd        enums.DeployCmd `in:"path" name:"cmd"`
}

func (r *ControlInstance) Path() string { return "/:instanceID/:cmd" }

func (r *ControlInstance) Output(ctx context.Context) (interface{}, error) {
	if _, err := validateByInstance(ctx, r.InstanceID); err != nil {
		return nil, err
	}

	return nil, deploy.ControlInstance(ctx, r.InstanceID, r.Cmd)
}
