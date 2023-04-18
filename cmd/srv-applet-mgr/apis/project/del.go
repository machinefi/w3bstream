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
	ctx, err := middleware.CurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.ProjectProviderFromContext(ctx))
	if err != nil {
		return nil, err
	}

	prj := types.MustProjectFromContext(ctx)

	if err = blockchain.RemoveMonitor(ctx, prj.Name); err != nil {
		return nil, err
	}

	return nil, project.Remove(ctx)
}
