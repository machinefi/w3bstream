package deploy

import (
	"context"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
	"github.com/machinefi/w3bstream/pkg/modules/projectoperator"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
	custommetrics "github.com/machinefi/w3bstream/pkg/types/wasm/metrics"
)

func WithInstanceRuntimeContext(parent context.Context) (context.Context, error) {
	d := types.MustMgrDBExecutorFromContext(parent)
	ins := types.MustInstanceFromContext(parent)
	app := types.MustAppletFromContext(parent)
	ctx := contextx.WithContextCompose(
		confid.WithSFIDGeneratorContext(confid.MustSFIDGeneratorFromContext(parent)),
		types.WithInstanceContext(ins),
		types.WithAppletContext(app),
		types.WithLoggerContext(types.MustLoggerFromContext(parent)),
		types.WithWasmDBEndpointContext(types.MustWasmDBEndpointFromContext(parent)),
		types.WithRedisEndpointContext(types.MustRedisEndpointFromContext(parent)),
		types.WithTaskWorkerContext(types.MustTaskWorkerFromContext(parent)),
		types.WithTaskBoardContext(types.MustTaskBoardFromContext(parent)),
		types.WithMqttBrokerContext(types.MustMqttBrokerFromContext(parent)),
		types.WithETHClientConfigContext(types.MustETHClientConfigFromContext(parent)),
	)(context.Background())

	prj := &models.Project{RelProject: models.RelProject{ProjectID: app.ProjectID}}
	if err := prj.FetchByProjectID(d); err != nil {
		return nil, err
	}

	res := &models.Resource{RelResource: models.RelResource{ResourceID: app.ResourceID}}
	if err := res.FetchByResourceID(d); err != nil {
		return nil, err
	}

	prjOp, err := projectoperator.GetByProject(parent, prj.ProjectID)
	if err != nil && err != status.ProjectOperatorNotFound {
		return nil, err
	}
	accOp, err := operator.ListByCond(parent, &operator.CondArgs{AccountID: prj.RelAccount.AccountID})
	if err != nil {
		return nil, err
	}
	chainClient := wasm.NewChainClient(parent, prj, accOp, prjOp)

	logger := types.MustLoggerFromContext(parent).WithValues(
		"@src", "wasm",
		"@prj", prj.Name,
		"@app", app.Name,
	)
	metrics := custommetrics.NewCustomMetric(prj.AccountID.String(), prj.Name)

	ctx = contextx.WithContextCompose(
		types.WithProjectContext(prj),
		types.WithResourceContext(res),
		wasm.WithRedisPrefixContext(prj.Name),
		wasm.WithChainClientContext(chainClient),
		wasm.WithLoggerContext(logger),
		wasm.WithCustomMetricsContext(metrics),
	)(ctx)

	configs, err := config.List(parent, &config.CondArgs{
		RelIDs: []types.SFID{prj.ProjectID, app.AppletID, res.ResourceID, ins.InstanceID},
	})

	if err != nil {
		return nil, err
	}
	for _, c := range configs {
		if err = wasm.InitConfiguration(ctx, c); err != nil {
			return nil, status.ConfigInitFailed.StatusErr().WithDesc(err.Error())
		}
		ctx = c.WithContext(ctx)
	}

	return ctx, nil
}
