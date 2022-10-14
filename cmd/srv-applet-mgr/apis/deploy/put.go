package deploy

import (
	"context"

	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/httptransport/httpx"

	"github.com/iotexproject/w3bstream/pkg/enums"
	"github.com/iotexproject/w3bstream/pkg/modules/deploy"
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
