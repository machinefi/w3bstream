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
	ProjectName string `name:"projectName" in:"path"`
}

func (r *GetProjectSchema) Path() string {
	return "/:projectName/" + enums.CONFIG_TYPE__PROJECT_SCHEMA.String()
}

func (r *GetProjectSchema) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)
	cfg, err := config.GetValue(
		ctx, prj.ProjectID,
		enums.CONFIG_TYPE__PROJECT_SCHEMA,
	)
	if err != nil {
		return nil, err
	}
	return cfg.(*wasm.Schema), nil
}

type GetProjectEnv struct {
	httpx.MethodGet
	ProjectName string `name:"projectName" in:"path"`
}

func (r *GetProjectEnv) Path() string {
	return "/:projectName/" + enums.CONFIG_TYPE__PROJECT_ENV.String()
}

func (r *GetProjectEnv) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)
	cfg, err := config.GetValue(
		ctx, prj.ProjectID,
		enums.CONFIG_TYPE__PROJECT_ENV,
	)
	if err != nil {
		return nil, err
	}
	return cfg.(*wasm.Env), nil
}
