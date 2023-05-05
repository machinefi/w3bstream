package applet

import (
	"context"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
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

func RemoveBySFID(ctx context.Context, id types.SFID) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Applet{RelApplet: models.RelApplet{AppletID: id}}

	return sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			if err := m.DeleteByAppletID(d); err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			return strategy.Remove(ctx, &strategy.CondArgs{
				AppletIDs: types.SFIDs{id},
			})
		},
		func(d sqlx.DBExecutor) error {
			return deploy.RemoveByAppletSFID(ctx, m.AppletID)
		},
	).Do()
}

func Remove(ctx context.Context, r *CondArgs) error {
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
			summary := statusx.ErrorFields{}
			for i := range lst {
				v := &lst[i]
				if err = RemoveBySFID(ctx, v.AppletID); err != nil {
					se := statusx.FromErr(err)
					summary = append(summary, &statusx.ErrorField{
						In:    v.AppletID.String(),
						Field: se.Key,
						Msg:   se.Desc,
					})
				}
			}
			if len(summary) > 0 {
				return status.BatchRemoveAppletFailed.StatusErr().
					AppendErrorFields(summary...)
			}
			return nil
		},
	).Do()
}

func List(ctx context.Context, r *ListReq) (*ListRsp, error) {
	var (
		d    = types.MustMgrDBExecutorFromContext(ctx)
		err  error
		app  = &models.Applet{}
		ret  = &ListRsp{}
		cond = r.Condition()
	)

	if ret.Data, err = app.List(d, cond, r.Addition()); err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	if ret.Total, err = app.Count(d, cond); err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return ret, nil
}

func ListDetail(ctx context.Context, r *ListReq) (*ListDetailRsp, error) {
	var (
		lst *ListRsp
		err error
		ins *models.Instance
		res *models.Resource
		ret = &ListDetailRsp{}
	)

	lst, err = List(ctx, r)
	if err != nil {
		return nil, err
	}
	ret = &ListDetailRsp{Total: lst.Total}

	for i := range lst.Data {
		app := &lst.Data[i]
		ins, _ = deploy.GetByAppletSFID(ctx, app.AppletID)
		res, _ = resource.GetBySFID(ctx, app.ResourceID)
		ret.Data = append(ret.Data, detail(app, ins, res))
	}
	return ret, nil
}

func Create(ctx context.Context, r *CreateReq) (*CreateRsp, error) {
	var (
		res *models.Resource
		acc = types.MustAccountFromContext(ctx)
		raw []byte
		err error
	)

	filename := r.WasmName
	if filename == "" {
		filename = r.AppletName + ".wasm"
	}
	res, raw, err = resource.Create(ctx, acc.AccountID, r.File, filename, r.WasmMd5)
	if err != nil {
		return nil, err
	}
	ctx = types.WithResource(ctx, res)

	var (
		idg = confid.MustNewSFIDGenerator()
		prj = types.MustProjectFromContext(ctx)
		app = &models.Applet{
			RelApplet:   models.RelApplet{AppletID: idg.MustGenSFID()},
			RelProject:  models.RelProject{ProjectID: prj.ProjectID},
			RelResource: models.RelResource{ResourceID: res.ResourceID},
			AppletInfo:  models.AppletInfo{Name: idg.MustGenSFID().String()},
		}
		sty []models.Strategy
		ins *models.Instance
	)

	err = sqlx.NewTasks(types.MustMgrDBExecutorFromContext(ctx)).With(
		func(d sqlx.DBExecutor) error {
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
			rb := &deploy.CreateReq{Cache: r.WasmCache}
			if r.Deploy == datatypes.TRUE {
				ins, err = deploy.UpsertByCode(ctx, rb, raw, enums.INSTANCE_STATE__STARTED)
			} else {
				ins, err = deploy.Create(ctx, rb)
			}
			return err
		},
	).Do()
	if err != nil {
		return nil, err
	}

	return &CreateRsp{
		Applet:     app,
		Instance:   ins,
		Resource:   res,
		Strategies: sty,
	}, nil
}

func Update(ctx context.Context, r *UpdateReq) (*UpdateRsp, error) {
	var (
		d   = types.MustMgrDBExecutorFromContext(ctx)
		app = types.MustAppletFromContext(ctx)
		ins *models.Instance // maybe not deployed
		res *models.Resource
		sty []models.Strategy
		raw []byte
		err error
	)

	// create resource if needed
	if r.File != nil {
		acc := types.MustAccountFromContext(ctx)
		filename, md5 := r.Info.WasmName, r.Info.WasmMd5
		if filename == "" {
			filename = r.AppletName + ".wasm"
		}
		res, raw, err = resource.Create(ctx, acc.AccountID, r.File, filename, md5)
	}

	err = sqlx.NewTasks(d).With(
		// update strategy
		func(d sqlx.DBExecutor) error {
			sty = r.BuildStrategies(ctx)
			if len(sty) == 0 {
				return nil
			}
			if err = strategy.Remove(ctx, &strategy.CondArgs{
				AppletIDs: types.SFIDs{app.AppletID},
			}); err != nil {
				return err
			}
			if err = strategy.BatchCreate(ctx, r.BuildStrategies(ctx)); err != nil {
				return err
			}
			return nil
		},
		// update and deploy instance
		func(d sqlx.DBExecutor) error {
			if r.File == nil {
				return nil // instance state will not be changed
			}
			if r.Info.Deploy != datatypes.TRUE {
				return nil
			}
			ins, err = deploy.GetByAppletSFID(ctx, app.AppletID)
			if err != nil {
				return err
			}
			var rb *deploy.CreateReq
			if r.Info.WasmCache != nil {
				rb = &deploy.CreateReq{Cache: r.Info.WasmCache}
			}
			ins, err = deploy.UpsertByCode(ctx, rb, raw, ins.State, ins.InstanceID)
			return err
		},
		// update applet info
		func(d sqlx.DBExecutor) error {
			if r.Info.AppletName != "" {
				app.Name = r.Info.AppletName
			}
			app.ResourceID = res.ResourceID
			if err = app.UpdateByAppletID(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.AppletNameConflict
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
	).Do()
	if err != nil {
		return nil, err
	}

	return &UpdateRsp{app, ins, res, sty}, nil
}
