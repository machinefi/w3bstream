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
	ctx = middleware.MustCurrentAccountFromContext(ctx).WithAccount(ctx)

	prefix, err := middleware.ProjectNameModifier(ctx)
	if err != nil {
		return nil, err
	}
	r.Name = prefix + r.Name

	rsp, err := project.Create(ctx, &r.CreateReq)
	if err != nil {
		return nil, err
	}
	rsp.Name, _ = middleware.ProjectNameForDisplay(rsp.Name)
	return rsp, nil
}
