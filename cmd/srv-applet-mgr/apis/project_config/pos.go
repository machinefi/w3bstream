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

type CreateProjectSchema struct {
	httpx.MethodPost
	ProjectName string `name:"projectName" in:"path"`
	wasm.Schema `in:"body"`
}

func (r *CreateProjectSchema) Path() string {
	return "/:projectName/" + enums.CONFIG_TYPE__PROJECT_SCHEMA.String()
}

func (r *CreateProjectSchema) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)
	return nil, config.CreateOrUpdateConfig(ctx, prj.ProjectID, &r.Schema)
}

type CreateOrUpdateProjectEnv struct {
	httpx.MethodPost
	ProjectName string `name:"projectName" in:"path"`
	wasm.Env    `in:"body"`
}

func (r *CreateOrUpdateProjectEnv) Path() string {
	return "/:projectName/" + enums.CONFIG_TYPE__PROJECT_ENV.String()
}

func (r *CreateOrUpdateProjectEnv) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)
	return nil, config.CreateOrUpdateConfig(ctx, prj.ProjectID, &r.Env)
}
