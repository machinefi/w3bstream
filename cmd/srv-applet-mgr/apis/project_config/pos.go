package project_config

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/modules/mq"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

// CreateProjectSchema create project database schema. support postgres
// TODO support schema migration and more sql dialect
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
	r.Schema.WithName(r.ProjectName)

	prj := types.MustProjectFromContext(ctx)
	return config.CreateConfig(ctx, prj.ProjectID, &r.Schema)
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
	return config.CreateOrUpdateConfig(ctx, prj.ProjectID, &r.Env)
}

// CreateProjectMqttBroker create project mqtt broker
// TODO support modification for broker
type CreateProjectMqttBroker struct {
	httpx.MethodPost
	ProjectName     string `name:"projectName" in:"path"`
	wasm.MqttBroker `in:"body"`
}

func (r *CreateProjectMqttBroker) Path() string {
	return "/:projectName/" + enums.CONFIG_TYPE__PROJECT_MQTT_BROKER.String()
}

func (r *CreateProjectMqttBroker) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)
	mq.StopChannel(prj.Name)
	return config.CreateConfig(
		types.WithMqttMsgHandler(
			ctx,
			func(cli mqtt.Client, msg mqtt.Message) {
				event.OnEventReceivedFromMqtt(ctx, msg)
			},
		),
		prj.ProjectID, &r.MqttBroker,
	)
}
