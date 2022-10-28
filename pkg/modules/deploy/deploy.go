package deploy

import (
	"context"

	confid "github.com/iotexproject/Bumblebee/conf/id"
	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"
	"github.com/pkg/errors"

	"github.com/iotexproject/w3bstream/pkg/enums"
	"github.com/iotexproject/w3bstream/pkg/errors/status"
	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/modules/vm"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type CreateInstanceRsp struct {
	InstanceID    types.SFID          `json:"instanceID"`
	InstanceState enums.InstanceState `json:"instanceState"`
}

func CreateInstance(ctx context.Context, path string, appletID types.SFID) (*CreateInstanceRsp, error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	m := &models.Instance{RelApplet: models.RelApplet{AppletID: appletID}}

	_, l = l.Start(ctx, "CreateInstance")
	defer l.End()

	exists, err := GetInstanceByAppletID(ctx, appletID)
	if err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "GetInstanceByAppletID")
	}
	for _, i := range exists {
		if err := ControlInstance(ctx, i.InstanceID, enums.DEPLOY_CMD__REMOVE); err != nil {
			l.Error(err)
			return nil, err
		}
	}

	m.InstanceID = idg.MustGenSFID()

	mApp := &models.Applet{RelApplet: models.RelApplet{AppletID: appletID}}
	if err := mApp.FetchByAppletID(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err)
	}
	mPrj := &models.Project{RelProject: models.RelProject{ProjectID: mApp.ProjectID}}
	if err := mPrj.FetchByProjectID(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err)
	}

	ctx = types.WithProject(ctx, mPrj)
	ctx = types.WithApplet(ctx, mApp)
	err = vm.NewInstance(ctx, path, m.InstanceID)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	m.State = enums.INSTANCE_STATE__CREATED
	m.Path = path

	if err = m.Create(d); err != nil {
		l.Error(err)
		_ = vm.DelInstance(ctx, m.InstanceID)
		return nil, err
	}
	l.WithValues("instance", m.InstanceID).Info("created")

	return &CreateInstanceRsp{
		InstanceID:    m.InstanceID,
		InstanceState: m.State,
	}, nil
}

func ControlInstance(ctx context.Context, instanceID types.SFID, cmd enums.DeployCmd) (err error) {
	var (
		d = types.MustDBExecutorFromContext(ctx)
		l = types.MustLoggerFromContext(ctx)
		m *models.Instance
	)

	_, l = l.Start(ctx, "ControlInstance")
	defer l.End()

	defer func() {
		l = l.WithValues("instance", instanceID, "cmd", cmd.String())
		if err != nil {
			l.Error(err)
		} else {
			l.Info("done")
		}
	}()

	if m, err = GetInstanceByInstanceID(ctx, instanceID); err != nil {
		l.Error(err)
		return err
	}

	switch cmd {
	case enums.DEPLOY_CMD__REMOVE:
		if err = vm.DelInstance(ctx, instanceID); err != nil {
			l.Error(err)
			return err
		}
		if err = m.DeleteByInstanceID(d); err != nil {
			l.Error(err)
			return status.CheckDatabaseError(err, "DeleteInstanceByInstanceID")
		}
		return nil
	case enums.DEPLOY_CMD__STOP:
		if err = vm.StopInstance(ctx, instanceID); err != nil {
			l.Error(err)
			return err
		}
		m.State = enums.INSTANCE_STATE__STOPPED
		if err = m.UpdateByInstanceID(d); err != nil {
			l.Error(err)
			return status.CheckDatabaseError(err, "UpdateInstanceByInstanceID")
		}
		return nil
	case enums.DEPLOY_CMD__START:
		if err = vm.StartInstance(ctx, instanceID); err != nil {
			l.Error(err)
			return err
		}
		m.State = enums.INSTANCE_STATE__STARTED
		if err = m.UpdateByInstanceID(d); err != nil {
			l.Error(err)
			return status.CheckDatabaseError(err, "UpdateInstanceByInstanceID")
		}
		return nil
	case enums.DEPLOY_CMD__RESTART:
		if err = vm.StopInstance(ctx, instanceID); err != nil {
			l.Error(err)
			return err
		}
		if err = vm.StartInstance(ctx, instanceID); err != nil {
			l.Error(err)
			return err
		}
		m.State = enums.INSTANCE_STATE__STARTED
		if err = m.UpdateByInstanceID(d); err != nil {
			l.Error(err)
			return status.CheckDatabaseError(err, "UpdateInstanceByInstanceID")
		}
		return nil
	default:
		return status.BadRequest.StatusErr().WithDesc("unknown deploy command")
	}
}

func GetInstanceByInstanceID(ctx context.Context, instanceID types.SFID) (*models.Instance, error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.Instance{RelInstance: models.RelInstance{InstanceID: instanceID}}

	_, l = l.Start(ctx, "GetInstanceByInstanceID")
	defer l.End()

	if err := m.FetchByInstanceID(d); err != nil {
		return nil, status.CheckDatabaseError(err, "FetchInstanceByInstanceID")
	}

	state, ok := vm.GetInstanceState(instanceID)
	if !ok {
		return nil, status.NotFound.StatusErr().WithDesc("instance not found in mgr")
	}
	if state != m.State {
		l.WithValues("mgr_state", state, "db_state", m.State).
			Warn(errors.New("unmatched"))
		m.State = state
		if err := m.UpdateByInstanceID(d); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func GetInstanceByAppletID(ctx context.Context, appletID types.SFID) (ret []models.Instance, err error) {
	d := types.MustDBExecutorFromContext(ctx)
	m := &models.Instance{}

	err = d.QueryAndScan(
		builder.Select(nil).From(
			d.T(m),
			builder.Where(m.ColAppletID().Eq(appletID)),
		),
		&ret,
	)
	return
}

func StartInstances(ctx context.Context) error {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.Instance{}

	_, l = l.Start(ctx, "StartInstances")
	defer l.End()

	list, err := m.List(d, nil)
	if err != nil {
		return err
	}
	for _, i := range list {
		l = l.WithValues("instance", i.InstanceID, "applet", i.AppletID)

		mApp := &models.Applet{RelApplet: models.RelApplet{AppletID: i.AppletID}}
		if err := mApp.FetchByAppletID(d); err != nil {
			l.Warn(err)
			continue
		}
		mPrj := &models.Project{RelProject: models.RelProject{ProjectID: mApp.ProjectID}}
		if err := mPrj.FetchByProjectID(d); err != nil {
			l.Warn(err)
			continue
		}

		ctx = types.WithProject(ctx, mPrj)
		ctx = types.WithApplet(ctx, mApp)
		err = vm.NewInstance(ctx, i.Path, i.InstanceID)
		cmd := enums.DEPLOY_CMD_UNKNOWN

		if err != nil {
			l.Warn(err)
			cmd = enums.DEPLOY_CMD__REMOVE
		} else {
			switch i.State {
			case enums.INSTANCE_STATE__CREATED:
				l.Info("created")
				continue
			case enums.INSTANCE_STATE__STARTED:
				cmd = enums.DEPLOY_CMD__START
			case enums.INSTANCE_STATE__STOPPED:
				cmd = enums.DEPLOY_CMD__STOP
			}
		}

		l = l.WithValues("cmd", cmd)
		if err = ControlInstance(ctx, i.InstanceID, cmd); err != nil {
			l.Error(err)
		}
	}
	return nil
}
