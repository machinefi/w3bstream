package project

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func CreateOrUpdateProjectEnv(ctx context.Context, env *wasm.Env) error {
	prj := types.MustProjectFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "CreateOrUpdateProjectEnv")

	val, err := json.Marshal(env)
	if err != nil {
		l.Error(err)
		return status.InternalServerError.StatusErr().WithDesc(err.Error())
	}

	_, err = config.CreateOrUpdateConfig(ctx, prj.ProjectID, enums.CONFIG_TYPE__PROJECT_ENV, val)
	return err
}

func CreateProjectSchema(ctx context.Context, schema *wasm.Schema) error {
	prj := types.MustProjectFromContext(ctx)
	l := types.MustLoggerFromContext(ctx).WithValues("project_id", prj.ProjectID)

	_, l = l.Start(ctx, "CreateProjectSchema")
	defer l.End()

	schema.WithName(prj.Name)
	if err := schema.Init(); err != nil {
		l.Error(err)
		return status.InternalServerError.StatusErr().
			WithDesc(fmt.Sprintf("init schema failed: [project:%s] [err:%v]",
				prj.Name, err))
	}
	if err := config.CreateConfig(ctx, prj.ProjectID, schema); err != nil {
		l.Error(err)
		return err
	}

	wasmdb := types.MustWasmDBExecutorFromContext(ctx)
	if _, err := wasmdb.Exec(schema.CreateSchema()); err != nil {
		l.Error(err)
		return status.InternalServerError.StatusErr().
			WithDesc(fmt.Sprintf("create wasm schema failed: [project:%s] [err:%v]",
				prj.Name, err))
	}

	db := schema.DBExecutor(wasmdb)
	for _, t := range schema.Tables {
		es := t.CreateIfNotExists()
		for _, e := range es {
			if e.IsNil() {
				continue
			}
			l.Info(builder.ResolveExpr(e).Query())
			if _, err := db.Exec(e); err != nil {
				l.Error(err)
				return status.InternalServerError.StatusErr().
					WithDesc(fmt.Sprintf("create wasm tables failed: [project:%s] [tbl:%s] [err:%v]",
						prj.Name, t.Name, err))
			}
		}
	}

	return nil
}
