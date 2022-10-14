package applet

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	confid "github.com/iotexproject/Bumblebee/conf/id"
	"github.com/iotexproject/Bumblebee/kit/sqlx"
	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"
	"github.com/iotexproject/Bumblebee/kit/sqlx/datatypes"

	"github.com/iotexproject/w3bstream/pkg/modules/vm"

	"github.com/iotexproject/w3bstream/pkg/errors/status"
	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/modules/resource"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type CreateAppletReq struct {
	File *multipart.FileHeader `name:"file"`
	Info `name:"info"`
}

type Info struct {
	AppletName string                `json:"appletName"`
	Strategies []models.StrategyInfo `json:"strategies,omitempty"`
}

func CreateApplet(ctx context.Context, projectID types.SFID, r *CreateAppletReq) (*models.Applet, error) {
	d := types.MustDBExecutorFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	appletID := idg.MustGenSFID()
	_, filename, sum, err := resource.Upload(ctx, r.File, appletID.String())
	if err != nil {
		return nil, status.UploadFileFailed.StatusErr().WithDesc(err.Error())
	}

	m := &models.Applet{
		RelProject: models.RelProject{ProjectID: projectID},
		RelApplet:  models.RelApplet{AppletID: appletID},
		AppletInfo: models.AppletInfo{
			Name: r.AppletName,
			Path: filename,
			Md5:  sum,
		},
	}
	if len(r.Info.Strategies) == 0 {
		r.Info.Strategies = append(r.Info.Strategies, models.DefaultStrategyInfo)
	}

	err = sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			return m.Create(db)
		},
		func(db sqlx.DBExecutor) error {
			for i := range r.Info.Strategies {
				if err := (&models.Strategy{
					RelStrategy:  models.RelStrategy{StrategyID: idg.MustGenSFID()},
					RelProject:   models.RelProject{ProjectID: projectID},
					RelApplet:    models.RelApplet{AppletID: m.AppletID},
					StrategyInfo: r.Info.Strategies[i],
				}).Create(db); err != nil {
					return err
				}
			}
			return nil
		},
	).Do()

	if err != nil {
		defer os.RemoveAll(filename)
		return nil, status.CheckDatabaseError(err, "CreateApplet")
	}

	return m, nil
}

type UpdateAppletReq struct {
	File  *multipart.FileHeader `name:"file"`
	*Info `name:"info"`
}

func UpdateApplet(ctx context.Context, appletID types.SFID, r *UpdateAppletReq) error {
	_, filename, sum, err := resource.Upload(ctx, r.File, appletID.String())
	if err != nil {
		return status.UploadFileFailed.StatusErr().WithDesc(err.Error())
	}

	d := types.MustDBExecutorFromContext(ctx)
	m := &models.Applet{RelApplet: models.RelApplet{AppletID: appletID}}
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	oldPath := ""
	needUpdateStrategies := r.Info != nil && len(r.Strategies) > 0

	sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			return m.FetchByAppletID(d)
		},
		func(db sqlx.DBExecutor) error {
			oldPath = m.Path
			m.Path, m.Md5 = filename, sum
			if r.Info != nil {
				m.Name = r.AppletName
			}
			return m.UpdateByAppletID(d)
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
							s.ColProjectID().Eq(m.ProjectID),
							s.ColAppletID().Eq(m.AppletID),
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
					RelProject:   models.RelProject{ProjectID: m.ProjectID},
					RelApplet:    models.RelApplet{AppletID: m.AppletID},
					StrategyInfo: r.Info.Strategies[i],
				}).Create(db); err != nil {
					return err
				}
			}
			return nil
		},
	).Do()

	if err != nil {
		os.RemoveAll(filename)
	} else {
		os.RemoveAll(oldPath)
	}

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

	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	l.Start(ctx, "ListApplets")
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
	ProjectID types.SFID `in:"path"  name:"projectID"`
	AppletID  types.SFID `in:"path"  name:"appletID"`
}

func RemoveApplet(ctx context.Context, r *RemoveAppletReq) error {
	var (
		d         = types.MustDBExecutorFromContext(ctx)
		mApplet   = &models.Applet{}
		mInstance = &models.Instance{}
		instances []models.Instance
		err       error
	)

	return sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			mApplet.AppletID = r.AppletID
			err = mApplet.FetchByAppletID(d)
			if err != nil {
				return status.CheckDatabaseError(err, "fetch by applet id")
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			mInstance.AppletID = r.AppletID
			instances, err = mInstance.List(d, mInstance.ColAppletID().Eq(r.AppletID))
			if err != nil {
				return status.CheckDatabaseError(err, "ListByAppletID")
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			for _, i := range instances {
				if err = vm.DelInstance(i.InstanceID.String()); err != nil {
					return status.InternalServerError.StatusErr().WithDesc(
						fmt.Sprintf("delete instance %s failed: %s",
							i.InstanceID, err.Error(),
						),
					)
				}
				if err = i.DeleteByInstanceID(d); err != nil {
					return status.CheckDatabaseError(err, "DeleteByInstanceID")
				}
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			return status.CheckDatabaseError(
				mApplet.DeleteByAppletID(d),
				"DeleteAppletByAppletID",
			)
		},
	).Do()
}

type GetAppletReq struct {
	ProjectID types.SFID `in:"path" name:"projectID"`
	AppletID  types.SFID `in:"path" name:"appletID"`
}

type GetAppletRsp struct {
	models.Applet
	Instances []models.Instance `json:"instances"`
}

func GetAppletByAppletID(ctx context.Context, appletID types.SFID) (*GetAppletRsp, error) {
	d := types.MustDBExecutorFromContext(ctx)
	m := &models.Applet{RelApplet: models.RelApplet{AppletID: appletID}}
	err := m.FetchByAppletID(d)
	if err != nil {
		return nil, status.CheckDatabaseError(err, "FetchByAppletID")
	}
	return &GetAppletRsp{Applet: *m}, err
}
