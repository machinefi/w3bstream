package monitor

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateMonitor struct {
	httpx.MethodPost
	ProjectID                   types.SFID `in:"path" name:"projectID"`
	blockchain.CreateMonitorReq `in:"body"`
}

func (r *CreateMonitor) Path() string { return "/:projectID" }

func (r *CreateMonitor) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	p, err := ca.ValidateProjectPerm(ctx, r.ProjectID)
	if err != nil {
		return nil, err
	}
	return blockchain.CreateMonitor(ctx, p.Name, &r.CreateMonitorReq)
}
