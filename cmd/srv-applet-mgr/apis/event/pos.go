package event

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/types"
)

type HandleEvent struct {
	httpx.MethodPost
	ProjectName          string `in:"path" name:"projectName"`
	event.HandleEventReq `in:"body"`
}

func (r *HandleEvent) Path() string { return "/:projectName" }

func (r *HandleEvent) Output(ctx context.Context) (interface{}, error) {
	var (
		err error
		ca  = middleware.CurrentAccountFromContext(ctx)
	)
	ctx, err = ca.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	return types.MustProjectFromContext(ctx), nil
}
