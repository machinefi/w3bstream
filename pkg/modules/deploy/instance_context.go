package deploy

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func WithInstanceRuntimeContext(parent context.Context) (context.Context, error) {
	d := types.MustMgrDBExecutorFromContext(parent)
	ins := types.MustInstanceFromContext(parent)
	ctx := contextx.WithContextCompose(
		types.WithInstanceContext(ins),
		types.WithLoggerContext(types.MustLoggerFromContext(parent)),
		types.WithWasmDBExecutorContext(types.MustWasmDBExecutorFromContext(parent)),
		types.WithRedisEndpointContext(types.MustRedisEndpointFromContext(parent)),
		types.WithMqttBrokerContext(types.MustMqttBrokerFromContext(parent)),
	)(context.Background())

	app := &models.Applet{RelApplet: models.RelApplet{AppletID: ins.AppletID}}
	if err := app.FetchByAppletID(d); err != nil {
		return nil, err
	}
	ctx = types.WithApplet(ctx, app)
	prj := &models.Project{RelProject: models.RelProject{ProjectID: app.ProjectID}}
	if err := prj.FetchByProjectID(d); err != nil {
		return nil, err
	}
	ctx = types.WithProject(ctx, prj)
	ctx = wasm.WithEnvPrefix(ctx, prj.Name)
	ctx = wasm.WithRedisPrefix(ctx, prj.Name)
	res := &models.Resource{RelResource: models.RelResource{ResourceID: app.ResourceID}}
	if err := res.FetchByResourceID(d); err != nil {
		return nil, err
	}
	ctx = types.WithResource(ctx, res)

	configs, err := config.FetchConfigValuesByRelIDs(parent, prj.ProjectID, app.AppletID, res.ResourceID, ins.InstanceID)
	if err != nil {
		return nil, err
	}
	for _, c := range configs {
		ctx = c.WithContext(ctx)
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
