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
	scm := &wasm.Schema{}
	if err = config.GetConfigValue(ctx, prj.ProjectID, scm); err != nil {
		return nil, err
	}
	return scm, nil
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
	env := &wasm.Env{}
	if err = config.GetConfigValue(ctx, prj.ProjectID, env); err != nil {
		return nil, err
	}
	return env, nil
}

type GetProjectMattBroker struct {
	httpx.MethodGet
	ProjectName string `name:"projectName" in:"path"`
}

func (r *GetProjectMattBroker) Path() string {
	return "/:projectName/" + enums.CONFIG_TYPE__PROJECT_MQTT_BROKER.String()
}

func (r *GetProjectMattBroker) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)
	broker := &wasm.MqttBroker{}
	if err = config.GetConfigValue(ctx, prj.ProjectID, broker); err != nil {
		return nil, err
	}
	return broker, nil
}

type GetProjectMattClient struct {
	httpx.MethodGet
	ProjectName string `name:"projectName" in:"path"`
}

func (r *GetProjectMattClient) Path() string {
	return "/:projectName/" + enums.CONFIG_TYPE__PROJECT_MQTT_CLIENT.String()
}

func (r *GetProjectMattClient) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)
	client := &wasm.MqttClient{}
	if err = config.GetConfigValue(ctx, prj.ProjectID, client); err != nil {
		return nil, err
	}
	return client, nil
}
