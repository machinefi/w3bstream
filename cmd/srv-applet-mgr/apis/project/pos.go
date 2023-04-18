package project

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/project"
)

type CreateProject struct {
	httpx.MethodPost
	project.CreateReq `in:"body"`
}

func (r *CreateProject) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	return project.CreateWithConfig(ctx, ca.AccountID, &r.CreateReq)
}

// func(ctx context.Context, channel string, data *eventpb.Event) (interface{}, error) {
// return event.OnEventReceived(ctx, channel, data)
// },
