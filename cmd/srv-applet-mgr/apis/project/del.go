package project

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/types"
)

type RemoveProject struct {
	httpx.MethodDelete
}

func (r *RemoveProject) Output(ctx context.Context) (interface{}, error) {
	name := middleware.MustProjectName(ctx)
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, name)
	if err != nil {
		return nil, err
	}
	// TODO @zhiran  move this to bff request
	if err := blockchain.RemoveMonitor(ctx, name); err != nil {
		return nil, err
	}

	v := types.MustProjectFromContext(ctx)
	return nil, project.RemoveBySFID(ctx, v.ProjectID)
}
