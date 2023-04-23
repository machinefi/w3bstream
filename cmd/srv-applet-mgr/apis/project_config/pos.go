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
	wasm.Schema `in:"body"`
}

func (r *CreateProjectSchema) Path() string {
	return "/" + enums.CONFIG_TYPE__PROJECT_SCHEMA.String()
}

func (r *CreateProjectSchema) Output(ctx context.Context) (interface{}, error) {
	prj := middleware.MustProjectName(ctx)
	ca := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, prj)
	if err != nil {
		return nil, err
	}
	return config.Create(ctx, types.MustProjectFromContext(ctx).ProjectID, &r.Schema)
}

type CreateOrUpdateProjectEnv struct {
	httpx.MethodPost
	wasm.Env `in:"body"`
}

func (r *CreateOrUpdateProjectEnv) Path() string {
	return "/" + enums.CONFIG_TYPE__PROJECT_ENV.String()
}

func (r *CreateOrUpdateProjectEnv) Output(ctx context.Context) (interface{}, error) {
	prj := middleware.MustProjectName(ctx)
	ca := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, prj)
	if err != nil {
		return nil, err
	}
	return config.Upsert(ctx, types.MustProjectFromContext(ctx).ProjectID, &r.Env)
}
