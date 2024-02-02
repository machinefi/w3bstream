package deploy

import (
	"context"
	"fmt"
	"strings"

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
	apisrv := types.MustWasmApiServerFromContext(parent)
	account := prj.AccountID.String()
	if strings.HasPrefix(prj.Name, "eth_") {
		parts := strings.Split(prj.Name, "_")
		if len(parts) >= 3 {
			account = strings.Join(parts[0:2], "_")
		}
	}
	metric := metrics.NewCustomMetric(account, prj.Name)
	logger := types.MustLoggerFromContext(parent)
	sfid := confid.MustSFIDGeneratorFromContext(parent)

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
		parent = c.WithContext(parent)
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
		parent = c.WithContext(parent)
	}
	chainconf := types.MustChainConfigFromContext(parent)
	operators := types.MustOperatorPoolFromContext(parent)

	return contextx.WithContextCompose(
		types.WithMgrDBExecutorContext(d),
		types.WithWasmApiServerContext(apisrv),
		types.WithLoggerContext(logger),
		wasm.WithCustomMetricsContext(metric),
		confid.WithSFIDGeneratorContext(sfid),
		types.WithProjectContext(prj),
		types.WithAppletContext(app),
		types.WithInstanceContext(ins),
		types.WithChainConfigContext(chainconf),
		types.WithOperatorPoolContext(operators),
	)(parent), nil
}
