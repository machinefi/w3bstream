package applet

import (
	"context"
	"fmt"
	"mime/multipart"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/machinefi/w3bstream/pkg/types"
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

func CreateApplet(ctx context.Context, projectID types.SFID, r *CreateAppletReq) (mApplet *models.Applet, err error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "CreateApplet")
	defer l.End()

	mResource := &models.Resource{}
	if mResource, err = resource.FetchOrCreateResource(ctx, r.File); err != nil {
		l.Error(err)
		return nil, err
	}

	appletID := idg.MustGenSFID()
	mApplet = &models.Applet{
		RelProject:  models.RelProject{ProjectID: projectID},
		RelApplet:   models.RelApplet{AppletID: appletID},
		RelResource: models.RelResource{ResourceID: mResource.RelResource.ResourceID},
		AppletInfo:  models.AppletInfo{Name: r.AppletName, WasmName: r.WasmName},
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

func UpdateApplet(ctx context.Context, appletID types.SFID, r *UpdateAppletReq) (err error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	mApplet := &models.Applet{RelApplet: models.RelApplet{AppletID: appletID}}
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "UpdateApplet")
	defer l.End()

	mResource := &models.Resource{}
	if mResource, err = resource.FetchOrCreateResource(ctx, r.File); err != nil {
		l.Error(err)
		return err
	}

	needUpdateStrategies := r.Info != nil && len(r.Strategies) > 0

	err = sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			return mApplet.FetchByAppletID(db)
		},
		func(db sqlx.DBExecutor) error {
			mApplet.RelResource = mResource.RelResource
			mApplet.WasmName = r.WasmName
			if r.Info != nil {
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

type ListAppletReq struct {
	ProjectName string       `in:"query" name:"projectName"`
	IDs         []uint64     `in:"query" name:"id,omitempty"`
	AppletIDs   []types.SFID `in:"query" name:"appletID,omitempty"`
	Names       []string     `in:"query" name:"names,omitempty"`
	NameLike    string       `in:"query" name:"name,omitempty"`
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

type RemoveAppletReq struct {
	AppletID types.SFID `in:"path"  name:"appletID"`
}

func RemoveApplet(ctx context.Context, r *RemoveAppletReq) error {
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
			mApplet.AppletID = r.AppletID
			err = mApplet.FetchByAppletID(d)
			if err != nil {
				l.Error(err)
				return status.CheckDatabaseError(err, "fetch by applet id")
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			mInstance.AppletID = r.AppletID
			instances, err = mInstance.List(d, mInstance.ColAppletID().Eq(r.AppletID))
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

type Detail struct {
	models.Applet
	models.ResourceInfo
	*models.InstanceInfo
}

func GetDetail(ctx context.Context) *Detail {
	app := types.MustAppletFromContext(ctx)
	res := types.MustResourceFromContext(ctx)
	ins, _ := types.InstanceFromContext(ctx)

	ret := &Detail{
		Applet:       *app,
		ResourceInfo: res.ResourceInfo,
	}

	if ins != nil {
		ret.InstanceInfo = &ins.InstanceInfo
	}

	return ret
}

func RemoveBySFID(ctx context.Context, id types.SFID) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Applet{RelApplet: models.RelApplet{AppletID: id}}

	if err := m.DeleteByAppletID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return status.AppletNotFound
		}
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func RemoveAppletAndInstanceBySFID(ctx context.Context, app *models.Applet, ins *models.Instance) error {
	d := types.MustMgrDBExecutorFromContext(ctx)

	if ins != nil {
		if err := vm.StopInstance(ctx, ins.InstanceID); err != nil {
			return status.StopInstanceFailed.StatusErr().WithDesc(err.Error())
		}
		if err := vm.DelInstance(ctx, ins.InstanceID); err != nil {
			return status.DeleteInstanceFailed.StatusErr().WithDesc(err.Error())
		}
	}

	return sqlx.NewTasks(d).With(
		func(_ sqlx.DBExecutor) error {
			return RemoveBySFID(ctx, app.AppletID)
		},
		func(_ sqlx.DBExecutor) error {
			if ins != nil {
				return deploy.RemoveBySFID(ctx, ins.InstanceID)
			}
			return nil
		},
	).Do()
}
