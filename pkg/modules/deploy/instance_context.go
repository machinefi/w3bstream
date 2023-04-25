package deploy

import (
	"context"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func WithInstanceRuntimeContext(parent context.Context) (context.Context, error) {
	d := types.MustMgrDBExecutorFromContext(parent)
	ins := types.MustInstanceFromContext(parent)
	ctx := contextx.WithContextCompose(
		confid.WithSFIDGeneratorContext(confid.MustSFIDGeneratorFromContext(parent)),
		types.WithInstanceContext(ins),
		types.WithLoggerContext(types.MustLoggerFromContext(parent)),
		types.WithWasmDBEndpointContext(types.MustWasmDBEndpointFromContext(parent)),
		types.WithRedisEndpointContext(types.MustRedisEndpointFromContext(parent)),
		types.WithTaskWorkerContext(types.MustTaskWorkerFromContext(parent)),
		types.WithTaskBoardContext(types.MustTaskBoardFromContext(parent)),
		types.WithMqttBrokerContext(types.MustMqttBrokerFromContext(parent)),
		types.WithETHClientConfigContext(types.MustETHClientConfigFromContext(parent)),
		types.WithFileSystemOpContext(types.MustFileSystemOpFromContext(parent)),
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
	ctx = wasm.WithRedisPrefix(ctx, prj.Name)
	res := &models.Resource{RelResource: models.RelResource{ResourceID: app.ResourceID}}
	if err := res.FetchByResourceID(d); err != nil {
		return nil, err
	}
	ctx = types.WithResource(ctx, res)

	configs, err := config.List(parent, &config.CondArgs{
		RelIDs: []types.SFID{prj.ProjectID, app.AppletID, res.ResourceID, ins.InstanceID},
	})

	if err != nil {
		return nil, err
	}
	for _, c := range configs {
		if canBeInit, ok := c.Configuration.(wasm.ConfigurationWithInit); ok {
			err = canBeInit.Init(ctx)
		}
		if err != nil {
			return nil, status.ConfigInitFailed.StatusErr().WithDesc(err.Error())
		}
		ctx = c.WithContext(ctx)
	}
	if _, ok := wasm.KVStoreFromContext(ctx); !ok {
		ctx = wasm.DefaultCache().WithContext(ctx)
	}
	if _, ok := wasm.MQTTClientFromContext(ctx); !ok {
		ctx = wasm.DefaultMQClient().WithContext(ctx)
	}

	acc := &models.Account{RelAccount: prj.RelAccount}
	if err := acc.FetchByAccountID(d); err != nil {
		return nil, err
	}
	ctx = wasm.WithChainClient(ctx, wasm.NewChainClient(ctx, &acc.OperatorPrivateKey))

	ctx = wasm.WithLogger(ctx, types.MustLoggerFromContext(ctx).WithValues(
		"@src", "wasm",
		"@prj", prj.Name,
		"@app", app.Name,
	))

	return ctx, nil
}
