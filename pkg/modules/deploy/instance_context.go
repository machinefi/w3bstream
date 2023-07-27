package deploy

import (
	"context"
	"fmt"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/modules/metrics"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
	"github.com/machinefi/w3bstream/pkg/modules/projectoperator"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func WithInstanceRuntimeContext(parent context.Context) (context.Context, error) {
	d := types.MustMgrDBExecutorFromContext(parent)
	ins := types.MustInstanceFromContext(parent)
	app := types.MustAppletFromContext(parent)

	var (
		prj    *models.Project
		exists bool
	)

	// completing parent context
	if prj, exists = types.ProjectFromContext(parent); !exists {
		prj = &models.Project{RelProject: models.RelProject{ProjectID: app.ProjectID}}
		if err := prj.FetchByProjectID(d); err != nil {
			return nil, err
		}
		parent = types.WithProject(parent, prj)
	}
	{
		op, err := projectoperator.GetByProject(parent, prj.ProjectID)
		if err != nil && err != status.ProjectOperatorNotFound {
			return nil, err
		}
		if op != nil {
			parent = types.WithProjectOperator(parent, op)
		}
		ops, err := operator.ListByCond(parent, &operator.CondArgs{AccountID: prj.RelAccount.AccountID})
		if err != nil {
			return nil, err
		}
		parent = types.WithOperators(parent, ops)
	}

	// wasm runtime context
	// original context is inherited from global context.
	// and configurations will be init from parent(host) context and with value to wasm runtime context
	ctx := contextx.WithContextCompose(
		types.WithWasmApiServerContext(types.MustWasmApiServerFromContext(parent)),
		types.WithLoggerContext(types.MustLoggerFromContext(parent)),
		types.WithTaskBoardContext(types.MustTaskBoardFromContext(parent)),
		types.WithTaskWorkerContext(types.MustTaskWorkerFromContext(parent)),
		confid.WithSFIDGeneratorContext(confid.MustSFIDGeneratorFromContext(parent)),
		wasm.WithCustomMetricsContext(metrics.NewCustomMetric(prj.AccountID.String(), prj.ProjectID.String())),
	)(context.Background())

	// with user defined contexts
	configs, err := config.List(parent, &config.CondArgs{
		RelIDs: []types.SFID{prj.ProjectID, app.AppletID, ins.InstanceID},
	})

	if err != nil {
		return nil, err
	}
	for _, c := range configs {
		if err = wasm.InitConfiguration(parent, c.Configuration); err != nil {
			return nil, status.ConfigInitFailed.StatusErr().WithDesc(
				fmt.Sprintf("config init failed: [type] %s [err] %v", c.ConfigType(), err),
			)
		}
		ctx = c.WithContext(ctx)
	}

	// with context from global configurations
	for _, t := range wasm.ConfigTypes {
		c, _ := wasm.NewGlobalConfigurationByType(t)
		if c == nil {
			continue
		}
		if err = wasm.InitGlobalConfiguration(parent, c); err != nil {
			return nil, status.ConfigInitFailed.StatusErr().WithDesc(
				fmt.Sprintf("global config init failed: [type] %s [err] %v", t, err),
			)
		}
		ctx = c.WithContext(ctx)
	}

	return ctx, nil
}
