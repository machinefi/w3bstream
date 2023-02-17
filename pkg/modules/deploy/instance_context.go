package deploy

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func WithInstanceRuntimeContext(parent context.Context) (context.Context, error) {
	d := types.MustMgrDBExecutorFromContext(parent)
	l := types.MustLoggerFromContext(parent)
	ins := types.MustInstanceFromContext(parent)
	ctx := contextx.WithContextCompose(
		types.WithInstanceContext(ins),
		types.WithLoggerContext(l),
		types.WithWasmDBExecutorContext(types.MustWasmDBExecutorFromContext(parent)),
		types.WithWasmPgEndpointContext(types.MustWasmPgEndpointFromContext(parent)),
		types.WithRedisEndpointContext(types.MustRedisEndpointFromContext(parent)),
		conflog.WithLoggerContext(l),
	)(context.Background())

	app, ok := types.AppletFromContext(parent)
	if !ok {
		app = &models.Applet{RelApplet: models.RelApplet{AppletID: ins.AppletID}}
		if err := app.FetchByAppletID(d); err != nil {
			return nil, err
		}
		parent = types.WithApplet(parent, app)
	}
	ctx = types.WithApplet(ctx, app)

	prj, ok := types.ProjectFromContext(parent)
	if !ok {
		prj = &models.Project{RelProject: models.RelProject{ProjectID: app.ProjectID}}
		if err := prj.FetchByProjectID(d); err != nil {
			return nil, err
		}
		parent = types.WithProject(parent, prj)
	}
	ctx = types.WithProject(ctx, prj)

	res, ok := types.ResourceFromContext(parent)
	if !ok {
		res = &models.Resource{RelResource: models.RelResource{ResourceID: app.ResourceID}}
		if err := res.FetchByResourceID(d); err != nil {
			return nil, err
		}
		parent = types.WithResource(parent, res)
	}
	ctx = types.WithResource(ctx, res)

	ctx = wasm.WithEnvPrefix(ctx, prj.Name)
	ctx = wasm.WithRedisPrefix(ctx, prj.Name)
	parent = types.WithMqttMsgHandler(parent, func(cli mqtt.Client, msg mqtt.Message) {
		event.OnEventReceivedFromMqtt(parent, msg)
	})

	configs, err := config.FetchConfigValuesByRelIDs(parent, prj.ProjectID, app.AppletID, res.ResourceID, ins.InstanceID)
	if err != nil {
		return nil, err
	}
	for _, c := range configs {
		ctx = c.WithContext(ctx)
		if err = wasm.Init(parent, c); err != nil {
			return nil, err
		}
	}
	if _, ok := wasm.KVStoreFromContext(ctx); !ok {
		ctx = wasm.DefaultCache().WithContext(ctx)
	}
	ctx = wasm.WithChainClient(ctx, wasm.NewChainClient(parent))
	ctx = wasm.WithLogger(ctx, types.MustLoggerFromContext(ctx).WithValues(
		"@src", "wasm",
		"@prj", prj.Name,
		"@app", app.Name,
	))
	return ctx, nil
}
