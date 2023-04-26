package applet

import (
	"context"
	"fmt"
	"mime/multipart"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

type CreateAppletReq struct {
	File *multipart.FileHeader `name:"file"`
	Info `name:"info"`
}

type Info struct {
	AppletName string                `json:"appletName"`
	WasmName   string                `json:"wasmName"`
	Strategies []models.StrategyInfo `json:"strategies,omitempty"`
}

type InfoApplet struct {
	models.Applet
	Path string `db:"f_wasm_path" json:"-"`
}

func CreateApplet(ctx context.Context, projectID, accountID types.SFID, r *CreateAppletReq) (mApplet *models.Applet, err error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "CreateApplet")
	defer l.End()

	mResource := &models.Resource{}
	if mResource, _, err = resource.Create(ctx, accountID, r.File, ""); err != nil {
		l.Error(err)
		return nil, err
	}

	appletID := idg.MustGenSFID()
	mApplet = &models.Applet{
		RelProject:  models.RelProject{ProjectID: projectID},
		RelApplet:   models.RelApplet{AppletID: appletID},
		RelResource: models.RelResource{ResourceID: mResource.RelResource.ResourceID},
		AppletInfo:  models.AppletInfo{Name: r.AppletName},
	}
	if len(r.Info.Strategies) == 0 {
		r.Info.Strategies = append(r.Info.Strategies, models.DefaultStrategyInfo)
	}

	err = sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			return mApplet.Create(db)
		},
		func(db sqlx.DBExecutor) error {
			for i := range r.Info.Strategies {
				if err := (&models.Strategy{
					RelStrategy:  models.RelStrategy{StrategyID: idg.MustGenSFID()},
					RelProject:   models.RelProject{ProjectID: projectID},
					RelApplet:    models.RelApplet{AppletID: mApplet.AppletID},
					StrategyInfo: r.Info.Strategies[i],
				}).Create(db); err != nil {
					return err
				}
			}
			return nil
		},
	).Do()

	if err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "CreateApplet")
	}

	return mApplet, nil
}

type UpdateAppletReq struct {
	File  *multipart.FileHeader `name:"file"`
	*Info `name:"info"`
}

func UpdateApplet(ctx context.Context, appletID, accountID types.SFID, r *UpdateAppletReq) (err error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	mApplet := &models.Applet{RelApplet: models.RelApplet{AppletID: appletID}}
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "UpdateApplet")
	defer l.End()

	mResource := &models.Resource{}
	if mResource, _, err = resource.Create(ctx, accountID, r.File, ""); err != nil {
		l.Error(err)
		return err
	}

	needUpdateAppletName := r.Info != nil && len(r.Info.AppletName) > 0
	needUpdateStrategies := r.Info != nil && len(r.Strategies) > 0

	err = sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			return mApplet.FetchByAppletID(db)
		},
		func(db sqlx.DBExecutor) error {
			mApplet.RelResource = mResource.RelResource
			if needUpdateAppletName {
				mApplet.Name = r.AppletName
			}
			return mApplet.UpdateByAppletID(db)
		},
		func(db sqlx.DBExecutor) error {
			if !needUpdateStrategies {
				return nil
			}
			s := &models.Strategy{}
			_, err := db.Exec(
				builder.Delete().From(
					db.T(s),
					builder.Where(
						builder.And(
							s.ColProjectID().Eq(mApplet.ProjectID),
							s.ColAppletID().Eq(mApplet.AppletID),
						),
					),
				),
			)
			return err
		},
		func(db sqlx.DBExecutor) error {
			for i := range r.Info.Strategies {
				if err := (&models.Strategy{
					RelStrategy:  models.RelStrategy{StrategyID: idg.MustGenSFID()},
					RelProject:   models.RelProject{ProjectID: mApplet.ProjectID},
					RelApplet:    models.RelApplet{AppletID: mApplet.AppletID},
					StrategyInfo: r.Info.Strategies[i],
				}).Create(db); err != nil {
					return err
				}
			}
			return nil
		},
	).Do()

	if err != nil {
		l.Error(err)
		return status.CheckDatabaseError(err, "UpdateApplet")
	}

	l.WithValues("applet", appletID, "path", mResource.ResourceInfo.Path).Info("applet uploaded")
	return nil
}

type UpdateAndDeployReq struct {
	File  *multipart.FileHeader `name:"file"`
	*Info `name:"info"`
	Cache *wasm.Cache `name:"cache,omitempty"`
}

func UpdateAndDeploy(ctx context.Context, accountID types.SFID, r *UpdateAndDeployReq) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)
	app := types.MustAppletFromContext(ctx)
	ins := types.MustInstanceFromContext(ctx)

	_, l = l.Start(ctx, "UpdateAndDeploy")
	defer l.End()

	mResource, code, err := resource.Create(ctx, accountID, r.File, "")
	if err != nil {
		l.Error(err)
		return err
	}

	needUpdateAppletName := r.Info != nil && len(r.Info.AppletName) > 0
	needUpdateStrategies := r.Info != nil && len(r.Info.Strategies) > 0

	_ctx := context.Background()
	err = sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			app.RelResource = mResource.RelResource
			if needUpdateAppletName {
				app.Name = r.AppletName
			}
			return app.UpdateByAppletID(db)
		},
		func(db sqlx.DBExecutor) error {
			if !needUpdateStrategies {
				return nil
			}
			s := &models.Strategy{}
			_, err := db.Exec(
				builder.Delete().From(
					db.T(s),
					builder.Where(
						builder.And(
							s.ColProjectID().Eq(app.ProjectID),
							s.ColAppletID().Eq(app.AppletID),
						),
					),
				),
			)
			return err
		},
		func(db sqlx.DBExecutor) error {
			for i := range r.Info.Strategies {
				if err := (&models.Strategy{
					RelStrategy:  models.RelStrategy{StrategyID: idg.MustGenSFID()},
					RelProject:   models.RelProject{ProjectID: app.ProjectID},
					RelApplet:    models.RelApplet{AppletID: app.AppletID},
					StrategyInfo: r.Info.Strategies[i],
				}).Create(db); err != nil {
					return err
				}
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			if r.Cache == nil {
				r.Cache = wasm.DefaultCache()
			}
			_, err = config.Create(ctx, ins.InstanceID, r.Cache)
			return err
		},
		func(db sqlx.DBExecutor) error {
			var _err error
			_ctx, _err = deploy.WithInstanceRuntimeContext(ctx)
			return _err
		},
		func(db sqlx.DBExecutor) error {
			return vm.NewInstance(_ctx, code, ins.InstanceID, enums.INSTANCE_STATE__STARTED)
		},
	).Do()

	if err != nil {
		l.Error(err)
		return status.CheckDatabaseError(err, "UpdateApplet")
	}

	l.WithValues("applet", app.AppletID, "path", mResource.ResourceInfo.Path).Info("applet uploaded and redeploy")
	return nil
}

type ListAppletReq struct {
	IDs       []uint64     `in:"query" name:"id,omitempty"`
	AppletIDs []types.SFID `in:"query" name:"appletID,omitempty"`
	Names     []string     `in:"query" name:"names,omitempty"`
	NameLike  string       `in:"query" name:"name,omitempty"`
	datatypes.Pager
}

func (r *ListAppletReq) Condition() builder.SqlCondition {
	var (
		m  = &models.Applet{}
		cs []builder.SqlCondition
	)
	if len(r.IDs) > 0 {
		cs = append(cs, m.ColID().In(r.IDs))
	}
	if len(r.AppletIDs) > 0 {
		cs = append(cs, m.ColAppletID().In(r.AppletIDs))
	}
	if len(r.Names) > 0 {
		cs = append(cs, m.ColName().In(r.Names))
	}
	if r.NameLike != "" {
		cs = append(cs, m.ColName().Like(r.NameLike))
	}
	return builder.And(cs...)
}

func (r *ListAppletReq) Additions() builder.Additions {
	m := &models.Applet{}
	return builder.Additions{
		builder.OrderBy(builder.DescOrder(m.ColCreatedAt())),
		r.Pager.Addition(),
	}
}

type ListAppletRsp struct {
	Data  []models.Applet `json:"data"`
	Hints int64           `json:"hints"`
}

func ListApplets(ctx context.Context, r *ListAppletReq) (*ListAppletRsp, error) {
	applet := &models.Applet{}

	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "ListApplets")
	defer l.End()

	applets, err := applet.List(d, r.Condition(), r.Additions()...)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	hints, err := applet.Count(d, r.Condition())
	if err != nil {
		l.Error(err)
		return nil, err
	}
	return &ListAppletRsp{applets, hints}, nil
}

func RemoveApplet(ctx context.Context, id types.SFID) error {
	var (
		d         = types.MustMgrDBExecutorFromContext(ctx)
		l         = types.MustLoggerFromContext(ctx)
		mApplet   = &models.Applet{}
		mInstance = &models.Instance{}
		instances []models.Instance
		err       error
	)

	_, l = l.Start(ctx, "RemoveApplet")
	defer l.End()

	return sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			mApplet.AppletID = id
			err = mApplet.FetchByAppletID(d)
			if err != nil {
				l.Error(err)
				return status.CheckDatabaseError(err, "fetch by applet id")
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			mInstance.AppletID = id
			instances, err = mInstance.List(d, mInstance.ColAppletID().Eq(id))
			if err != nil {
				l.Error(err)
				return status.CheckDatabaseError(err, "ListByAppletID")
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			for _, i := range instances {
				if err = vm.DelInstance(ctx, i.InstanceID); err != nil {
					l.Error(err)
					return status.InternalServerError.StatusErr().WithDesc(
						fmt.Sprintf("delete instance %s failed: %s",
							i.InstanceID, err.Error(),
						),
					)
				}
				if err = i.DeleteByInstanceID(d); err != nil {
					l.Error(err)
					return status.CheckDatabaseError(err, "DeleteByInstanceID")
				}
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			err = mApplet.DeleteByAppletID(d)
			if err != nil {
				l.Error(err)
				return status.CheckDatabaseError(err, "DeleteAppletByAppletID")
			}
			return nil
		},
	).Do()
}

type GetAppletReq struct {
	ProjectName string     `in:"path" name:"projectName"`
	AppletID    types.SFID `in:"path" name:"appletID"`
}

type GetAppletRsp struct {
	InfoApplet
	Instances []models.Instance `json:"instances"`
}

func GetAppletByAppletID(ctx context.Context, appletID types.SFID) (*GetAppletRsp, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	mApplet := &models.Applet{RelApplet: models.RelApplet{AppletID: appletID}}
	mResource := &models.Resource{}

	_, l = l.Start(ctx, "GetAppletByAppletID")
	defer l.End()

	err := sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			err := mApplet.FetchByAppletID(d)
			if err != nil {
				l.Error(err)
				return status.CheckDatabaseError(err, "FetchByAppletID")
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			mResource.ResourceID = mApplet.ResourceID
			err := mResource.FetchByResourceID(d)
			if err != nil {
				l.Error(err)
				return status.CheckDatabaseError(err, "FetchByWasmResourceID")
			}
			return nil
		},
	).Do()

	if err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "GetAppletByAppletID")
	}

	return &GetAppletRsp{
		InfoApplet: InfoApplet{Applet: *mApplet,
			Path: mResource.ResourceInfo.Path,
		},
	}, err
}

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
