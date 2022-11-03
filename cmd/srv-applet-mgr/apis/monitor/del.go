package monitor

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
	"github.com/machinefi/w3bstream/pkg/types"
)

type RemoveMonitor struct {
	httpx.MethodDelete
	ProjectID                   types.SFID `in:"path" name:"projectID"`
	blockchain.RemoveMonitorReq `in:"body"`
}

func (r *RemoveMonitor) Path() string { return "/:projectID" }

func (r *RemoveMonitor) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	p, err := ca.ValidateProjectPerm(ctx, r.ProjectID)
	if err != nil {
		return nil, err
	}
	return nil, blockchain.RemoveMonitor(ctx, p.Name, &r.RemoveMonitorReq)
}
