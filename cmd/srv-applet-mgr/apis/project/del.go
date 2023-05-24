package project

import (
	"context"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
	"github.com/machinefi/w3bstream/pkg/modules/metrics"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/types"
)

type RemoveProject struct {
	httpx.MethodDelete
}

func (r *RemoveProject) Output(ctx context.Context) (interface{}, error) {
	name := middleware.MustProjectName(ctx)
	acc := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := acc.WithProjectContextByName(ctx, name)
	if err != nil {
		return nil, err
	}
	// TODO @zhiran  move this to bff request
	var errRM error
	if err := blockchain.RemoveMonitor(ctx, name); err != nil {
		errRM = errors.Wrap(err, errRM.Error())
		l := types.MustLoggerFromContext(ctx)
		l.Error(err)
	}

	metrics.RemoveMetrics(ctx, acc.AccountID.String(), name)

	if err := project.RemoveBySFID(ctx, types.MustProjectFromContext(ctx).ProjectID); err != nil {
		errRM = errors.Wrap(err, errRM.Error())
		l := types.MustLoggerFromContext(ctx)
		l.Error(err)
	}

	return nil, errRM
}
