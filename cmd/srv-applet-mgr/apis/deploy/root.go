package deploy

import (
	"context"

	"github.com/iotexproject/Bumblebee/kit/httptransport"
	"github.com/iotexproject/Bumblebee/kit/kit"
	"github.com/iotexproject/Bumblebee/kit/sqlx"
	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"

	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/iotexproject/w3bstream/pkg/errors/status"
	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/types"
)

var Root = kit.NewRouter(httptransport.Group("/deploy"))

func init() {
	Root.Register(kit.NewRouter(&CreateInstance{}))
	Root.Register(kit.NewRouter(&GetInstanceByInstanceID{}))
	Root.Register(kit.NewRouter(&GetInstanceByAppletID{}))
	Root.Register(kit.NewRouter(&ControlInstance{}))
}

func validateByInstance(ctx context.Context, instanceID types.SFID) (*models.Instance, error) {
	d := types.MustDBExecutorFromContext(ctx)
	ca := middleware.CurrentAccountFromContext(ctx)

	mInstance := &models.Instance{RelInstance: models.RelInstance{InstanceID: instanceID}}
	err := mInstance.FetchByInstanceID(d)
	if err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.NotFound.StatusErr().WithDesc("instance not found")
		}
		return nil, err
	}

	mApplet := &models.Applet{}
	tInstance := d.T(mInstance)
	tApplet := d.T(mApplet)

	mProject := make([]struct {
		ProjectID types.SFID `db:"t_applet_f_project_id"`
	}, 0)

	err = d.QueryAndScan(
		builder.Select(
			builder.MultiWith(",",
				builder.Alias(mApplet.ColProjectID(), "t_applet_f_project_id"),
			),
		).
			From(
				tInstance,
				builder.LeftJoin(tApplet).On(
					mInstance.ColAppletID().Eq(mApplet.ColAppletID()),
				),
				builder.Where(mInstance.ColInstanceID().Eq(instanceID)),
			),
		&mProject,
	)
	if err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.NotFound.StatusErr().WithDesc("project not found")
		}
		return nil, err
	}
	if len(mProject) == 0 {
		return nil, status.NotFound.StatusErr().WithDesc("project not found")
	}

	if _, err = ca.ValidateProjectPerm(ctx, mProject[0].ProjectID); err != nil {
		return nil, err
	}
	return mInstance, nil
}
