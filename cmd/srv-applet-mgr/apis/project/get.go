package project

import (
	"context"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/types"
)

type GetProject struct {
	httpx.MethodGet
}

func (r *GetProject) Path() string { return "/data" }

func (r *GetProject) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)
	prj.Name, err = middleware.ProjectNameForDisplay(prj.Name)
	if err != nil {
		return nil, status.DeprecatedProject.StatusErr().
			WithDesc("this project is deprecated")
	}
	return prj, nil
}

type ListProject struct {
	httpx.MethodGet
	project.ListReq
}

func (r *ListProject) Path() string { return "/datalist" }

func (r *ListProject) Output(ctx context.Context) (interface{}, error) {
	ctx = middleware.MustCurrentAccountFromContext(ctx).WithAccount(ctx)
	rsp, err := project.List(ctx, &r.ListReq)
	if err != nil {
		return nil, err
	}

	_, l := types.MustLoggerFromContext(ctx).Start(ctx, "ListProject")
	for i := 0; i < len(rsp.Data); i++ {
		v := &rsp.Data[i]
		v.Name, err = middleware.ProjectNameForDisplay(v.Name)
		if err != nil {
			l.WithValues("prj", v.ProjectID).Warn(
				errors.New("this project is deprecated, no prefix"),
			)
			rsp.Data = append(rsp.Data[0:i], rsp.Data[i+1:]...)
			continue
		}
	}
	return rsp, nil
}
