package project

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
	"github.com/machinefi/w3bstream/pkg/modules/metrics"
	"github.com/machinefi/w3bstream/pkg/modules/project"
)

type CreateProject struct {
	httpx.MethodPost
	project.CreateReq `in:"body"`
}

func (r *CreateProject) Output(ctx context.Context) (interface{}, error) {
	acc := middleware.MustCurrentAccountFromContext(ctx)
	ctx = acc.WithAccount(ctx)

	prefix, err := middleware.ProjectNameModifier(ctx)
	if err != nil {
		return nil, err
	}
	r.Name = prefix + r.Name

	// make sure monitor and metrics deleted
	blockchain.RemoveMonitor(ctx, r.Name)
	metrics.RemoveMetrics(ctx, acc.AccountID.String(), r.Name)

	rsp, err := project.Create(ctx, &r.CreateReq)
	if err != nil {
		return nil, err
	}
	rsp.Name, _ = middleware.ProjectNameForDisplay(rsp.Name)
	return rsp, nil
}
