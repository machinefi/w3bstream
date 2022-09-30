package applet

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
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
	ProjectID  string               `json:"projectID"`
	AppletName string               `json:"appletName"`
	Config     *models.AppletConfig `json:"config,omitempty"`
}

func CreateApplet(ctx context.Context, r *CreateAppletReq) (*models.Applet, error) {
	appletID := uuid.New().String()
	_, filename, err := resource.Upload(ctx, r.File, appletID)
	if err != nil {
		return nil, status.UploadFileFailed.StatusErr().WithDesc(err.Error())
	}

	d := types.MustDBExecutorFromContext(ctx)

	m := &models.Applet{
		RelProject: models.RelProject{ProjectID: r.ProjectID},
		RelApplet:  models.RelApplet{AppletID: appletID},
		AppletInfo: models.AppletInfo{Name: r.AppletName, Path: filename, Config: r.Config},
	}

	if err = m.Create(d); err != nil {
		defer os.RemoveAll(filename)
		return nil, err
	}

	return m, nil
}

type UpdateAppletReq struct {
	File *multipart.FileHeader `name:"file"`
}

func UpdateApplet(ctx context.Context, appletID string, r *UpdateAppletReq) error {
	_, filename, err := resource.Upload(ctx, r.File, appletID)
	if err != nil {
		return status.UploadFileFailed.StatusErr().WithDesc(err.Error())
	}

	d := types.MustDBExecutorFromContext(ctx)
	m := &models.Applet{RelApplet: models.RelApplet{AppletID: appletID}}

	if err := m.FetchByAppletID(d); err != nil {
		return status.CheckDatabaseError(err, "FetchAppletByAppletID")
	}
	m.Path = filename
	if err := m.UpdateByAppletID(d); err != nil {
		defer os.RemoveAll(filename)
		return status.CheckDatabaseError(err, "UpdateAppletByAppletID")
	}
	return nil
}

type ListAppletReq struct {
	IDs       []string `in:"query" name:"id,omitempty"`
	AppletIDs []string `in:"query" name:"appletID,omitempty"`
	Names     []string `in:"query" name:"names,omitempty"`
	NameLike  string   `in:"query" name:"name,omitempty"`
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
	ProjectID string `in:"path"  name:"projectID"`
	AppletID  string `in:"path"  name:"appletID"`
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
				if err = vm.DelInstance(i.InstanceID); err != nil {
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
				"delete applet by applet id",
			)
		},
	).Do()
}

type GetAppletReq struct {
	ProjectID string `in:"path" name:"projectID"`
	AppletID  string `in:"path" name:"appletID"`
}

type GetAppletRsp struct {
	models.Applet
	Instances []models.Instance `json:"instances"`
}

func GetAppletByAppletID(ctx context.Context, appletID string) (*GetAppletRsp, error) {
	d := types.MustDBExecutorFromContext(ctx)
	m := &models.Applet{RelApplet: models.RelApplet{AppletID: appletID}}
	err := m.FetchByAppletID(d)
	if err != nil {
		return nil, status.CheckDatabaseError(err, "FetchByAppletID")
	}
	return &GetAppletRsp{
		Applet: *m,
	}, err
}
