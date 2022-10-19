package monitor

import (
	"context"

	"github.com/iotexproject/Bumblebee/kit/httptransport/httpx"

	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/iotexproject/w3bstream/pkg/modules/project"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type CreateMonitor struct {
	httpx.MethodPost
	ProjectID                types.SFID `in:"path" name:"projectID"`
	project.CreateMonitorReq `in:"body"`
}

func (r *CreateMonitor) Path() string { return "/:projectID" }

func (r *CreateMonitor) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	p, err := ca.ValidateProjectPerm(ctx, r.ProjectID)
	if err != nil {
		return nil, err
	}
	return project.CreateMonitor(ctx, p.Name, &r.CreateMonitorReq)
}
