package applet

import (
	"context"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func GetBySFID(ctx context.Context, id types.SFID) (*models.Applet, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Applet{RelApplet: models.RelApplet{AppletID: id}}

	if err := m.FetchByAppletID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.AppletNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func GetDetail(_ context.Context, app *models.Applet, ins *models.Instance, res *models.Resource) *Detail {
	ret := &Detail{
		Applet:       *app,
		ResourceInfo: res.ResourceInfo,
	}

	if ins != nil {
		ins.State, _ = vm.GetInstanceState(ins.InstanceID)
		ret.InstanceInfo = &ins.InstanceInfo
	}

	return ret
}

func Remove(ctx context.Context) error {
	var (
		d      = types.MustMgrDBExecutorFromContext(ctx)
		app    = types.MustAppletFromContext(ctx)
		ins, _ = types.InstanceFromContext(ctx)
		err    error
	)

	return sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			if err := app.DeleteByAppletID(d); err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			if ins == nil {
				ins, err = deploy.GetByAppletSFID(ctx, app.AppletID)
				return err
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			if ins != nil {
				ctx = types.WithInstance(ctx, ins)
				return deploy.Remove(ctx, false)
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			return strategy.BatchRemove(ctx, &strategy.DataListParam{
				ListReq: strategy.ListReq{
					AppletIDs: []types.SFID{app.AppletID},
				},
			})
		},
	).Do()
}

func BatchRemove(ctx context.Context, r *DataListParam) error {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.Applet{}

		err error
		lst []models.Applet
	)

	return sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			lst, err = m.List(d, r.Condition())
			if err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			for i := range lst {
				ctx = types.WithApplet(ctx, &lst[i])
				if err = Remove(ctx); err != nil {
					return err
				}
			}
			return nil
		},
	).Do()
}

func List(ctx context.Context, prj types.SFID, r *ListReq) (*ListRsp, error) {
	var (
		d    = types.MustMgrDBExecutorFromContext(ctx)
		err  error
		app  = &models.Applet{}
		ret  = &ListRsp{}
		cond = r.Condition(prj)
		adds = r.Additions()
	)

	if ret.Data, err = app.List(d, cond, adds...); err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	if ret.Hints, err = app.Count(d, cond); err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return ret, nil
}

func ListDetail(ctx context.Context, prj types.SFID, r *ListReq) (*ListDetailRsp, error) {
	var (
		lst *ListRsp
		err error
		ins *models.Instance
		res *models.Resource
		ret = &ListDetailRsp{}
	)

	lst, err = List(ctx, prj, r)
	if err != nil {
		return nil, err
	}
	ret = &ListDetailRsp{Hints: lst.Hints}

	for i := range lst.Data {
		app := &lst.Data[i]
		if ins, err = deploy.GetByAppletSFID(ctx, app.AppletID); err != nil {
			if se, ok := statusx.IsStatusErr(err); !ok && !se.Is(status.InstanceNotFound) {
				return nil, err
			}
		}
		if res, err = resource.GetBySFID(ctx, app.ResourceID); err != nil {
			return nil, err
		}
		ret.Data = append(ret.Data, GetDetail(ctx, app, ins, res))
	}
	return ret, nil
}

func Create(ctx context.Context, acc types.SFID, r *CreateReq) (*CreateRsp, error) {
	var (
		res *models.Resource
		err error
	)

	res, err = resource.Create(ctx, acc, r.WasmMd5, r.File)
	if err != nil {
		return nil, err
	}
	ctx = types.WithResource(ctx, res)

	var (
		idg = confid.MustNewSFIDGenerator()
		prj = types.MustProjectFromContext(ctx)
		app *models.Applet
		ins *models.Instance
	)

	err = sqlx.NewTasks(types.MustMgrDBExecutorFromContext(ctx)).With(
		func(d sqlx.DBExecutor) error {
			app = &models.Applet{
				RelApplet:   models.RelApplet{AppletID: idg.MustGenSFID()},
				RelProject:  models.RelProject{ProjectID: prj.ProjectID},
				RelResource: models.RelResource{ResourceID: res.ResourceID},
				AppletInfo: models.AppletInfo{
					Name:     r.AppletName,
					WasmName: r.WasmName,
				},
			}
			if err = app.Create(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.AppletNameConflict
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			ctx = types.WithApplet(ctx, app)
			return nil
		},
		func(d sqlx.DBExecutor) error {
			return strategy.BatchCreate(ctx, r.BuildStrategies(ctx))
		},
		func(d sqlx.DBExecutor) error {
			if r.WasmCache == nil {
				r.WasmCache = wasm.DefaultCache()
			}
			ins, err = deploy.Create(ctx, &deploy.CreateReq{
				CacheConfig: r.WasmCache,
			})
			return err
		},
	).Do()
	if err != nil {
		return nil, err
	}

	return &CreateRsp{
		RelApplet:    app.RelApplet,
		AppletInfo:   app.AppletInfo,
		RelInstance:  ins.RelInstance,
		InstanceInfo: ins.InstanceInfo,
	}, nil
}

func Update(ctx context.Context, acc types.SFID, r *UpdateReq) (*UpdateRsp, error) {
	var (
		d   = types.MustMgrDBExecutorFromContext(ctx)
		app = types.MustAppletFromContext(ctx)
		res = types.MustResourceFromContext(ctx)
		ins *models.Instance // maybe not deployed
		err error

		keepInstanceConfig = false
		keepStrategies     = false
	)

	if r.Info == nil || len(r.Info.Strategies) == 0 {
		keepStrategies = true
	}
	if r.Info == nil || r.Info.WasmCache == nil {
		keepInstanceConfig = true
	}

	err = sqlx.NewTasks(d).With(
		// create resource
		func(d sqlx.DBExecutor) error {
			if r.File != nil {
				md5 := ""
				if r.Info != nil {
					md5 = r.Info.WasmMd5
				}
				res, err = resource.Create(ctx, acc, md5, r.File)
				if err != nil {
					return err
				}
				ctx = types.WithResource(ctx, res)
			}
			return err
		},
		// drop old strategies
		func(d sqlx.DBExecutor) error {
			if keepStrategies {
				return nil
			}
			return strategy.BatchRemove(ctx, &strategy.DataListParam{
				ListReq: strategy.ListReq{AppletIDs: []types.SFID{app.AppletID}},
			})
		},
		// create new strategies
		func(d sqlx.DBExecutor) error {
			if keepStrategies {
				return nil
			}
			return strategy.BatchCreate(ctx, r.BuildStrategies(ctx))
		},
		// update applet
		func(d sqlx.DBExecutor) error {
			if ass := r.Assignments(); len(ass) != 0 {
				_, err = d.Exec(builder.Update(d.T(app)).Set(ass...))
				return err
			}
			return nil
		},
		// fetch new applet and with context
		func(d sqlx.DBExecutor) error {
			if len(r.Assignments()) == 0 {
				return nil
			}
			app, err = GetBySFID(ctx, app.AppletID)
			if err != nil {
				return err
			}
			ctx = types.WithApplet(ctx, app)
			return nil
		},
		// remove old instance(for updating) and deploy
		func(d sqlx.DBExecutor) error {
			return deploy.Remove(ctx, keepInstanceConfig)
		},
		// create and deploy new instance
		func(d sqlx.DBExecutor) error {
			req := &deploy.CreateReq{}
			if r.Info != nil {
				req.CacheConfig = r.Info.WasmCache
			}
			ins, err = deploy.Create(ctx, req)
			return err
		},
	).Do()
	if err != nil {
		return nil, err
	}

	return &UpdateRsp{
		RelApplet:    app.RelApplet,
		AppletInfo:   app.AppletInfo,
		RelInstance:  ins.RelInstance,
		InstanceInfo: ins.InstanceInfo,
	}, nil
}
