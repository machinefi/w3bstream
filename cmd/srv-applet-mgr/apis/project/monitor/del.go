package monitor

import (
	"context"

	"github.com/iotexproject/Bumblebee/kit/httptransport/httpx"

	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/iotexproject/w3bstream/pkg/modules/project"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type RemoveMonitor struct {
	httpx.MethodDelete
	ProjectID                types.SFID `in:"path" name:"projectID"`
	project.RemoveMonitorReq `in:"body"`
}

func (r *RemoveMonitor) Path() string { return "/:projectID" }

func (r *RemoveMonitor) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	p, err := ca.ValidateProjectPerm(ctx, r.ProjectID)
	if err != nil {
		return nil, err
	}
	return nil, project.RemoveMonitor(ctx, p.Name, &r.RemoveMonitorReq)
}
