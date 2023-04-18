package project_config

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

type GetProjectSchema struct {
	httpx.MethodGet
}

func (r *GetProjectSchema) Path() string {
	return "/" + enums.CONFIG_TYPE__PROJECT_SCHEMA.String()
}

func (r *GetProjectSchema) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)
	scm := &wasm.Schema{}
	if err = config.GetConfigValue(ctx, prj.ProjectID, scm); err != nil {
		return nil, err
	}
	return scm, nil
}

type GetProjectEnv struct {
	httpx.MethodGet
}

func (r *GetProjectEnv) Path() string {
	return "/" + enums.CONFIG_TYPE__PROJECT_ENV.String()
}

func (r *GetProjectEnv) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)
	env := &wasm.Env{}
	if err = config.GetConfigValue(ctx, prj.ProjectID, env); err != nil {
		return nil, err
	}
	return env, nil
}
