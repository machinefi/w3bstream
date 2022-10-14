package deploy

import (
	"context"
	"fmt"

	confid "github.com/iotexproject/Bumblebee/conf/id"
	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"
	"github.com/pkg/errors"

	"github.com/iotexproject/w3bstream/pkg/enums"
	"github.com/iotexproject/w3bstream/pkg/errors/status"
	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/modules/vm"
	"github.com/iotexproject/w3bstream/pkg/modules/vm/common"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type CreateInstanceRsp struct {
	InstanceID    types.SFID          `json:"instanceID"`
	InstanceState enums.InstanceState `json:"instanceState"`
}

func CreateInstance(ctx context.Context, path string, appletID types.SFID) (*CreateInstanceRsp, error) {
	d := types.MustDBExecutorFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)
	m := &models.Instance{
		RelApplet: models.RelApplet{AppletID: appletID},
	}

	exists, err := GetInstanceByAppletID(ctx, appletID)
	if err != nil {
		return nil, status.CheckDatabaseError(err, "GetInstanceByAppletID")
	}
	for _, i := range exists {
		if err := ControlInstance(ctx, i.InstanceID, enums.DEPLOY_CMD__REMOVE); err != nil {
			return nil, err
		}
	}

	m.InstanceID = idg.MustGenSFID()

	err = vm.NewInstanceWithID(path, m.InstanceID.String(), common.DefaultInstanceOptionSetter)
	if err != nil {
		return nil, err
	}
	m.State = enums.INSTANCE_STATE__CREATED
	m.Path = path

	if err = m.Create(d); err != nil {
		_ = vm.DelInstance(m.InstanceID.String())
		return nil, err
	}

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

	l.Start(ctx, "ControlInstance")
	defer l.End()

	defer func() {
		if err != nil {
			l.WithValues("instance", instanceID, "cmd", cmd).Error(err)
		}
	}()

	if m, err = GetInstanceByInstanceID(ctx, instanceID); err != nil {
		return err
	}

	switch cmd {
	case enums.DEPLOY_CMD__REMOVE:
		if err = vm.DelInstance(instanceID.String()); err != nil {
			return err
		}
		return status.CheckDatabaseError(m.DeleteByInstanceID(d), "DeleteInstanceByInstanceID")
	case enums.DEPLOY_CMD__STOP:
		if err = vm.StopInstance(instanceID.String()); err != nil {
			return err
		}
		m.State = enums.INSTANCE_STATE__STOPPED
		return status.CheckDatabaseError(m.UpdateByInstanceID(d), "UpdateInstanceByInstanceID")
	case enums.DEPLOY_CMD__START:
		if err = vm.StartInstance(instanceID.String()); err != nil {
			return err
		}
		m.State = enums.INSTANCE_STATE__STARTED
		return status.CheckDatabaseError(m.UpdateByInstanceID(d), "UpdateInstanceByInstanceID")
	case enums.DEPLOY_CMD__RESTART:
		if err = vm.StopInstance(instanceID.String()); err != nil {
			return err
		}
		if err = vm.StartInstance(instanceID.String()); err != nil {
			return err
		}
		m.State = enums.INSTANCE_STATE__CREATED
		return status.CheckDatabaseError(m.UpdateByInstanceID(d), "UpdateInstanceByInstanceID")
	default:
		return status.BadRequest.StatusErr().WithDesc("unknown deploy command")
	}
}

func GetInstanceByInstanceID(ctx context.Context, instanceID types.SFID) (*models.Instance, error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.Instance{RelInstance: models.RelInstance{InstanceID: instanceID}}

	l.Start(ctx, "GetInstanceByInstanceID")

	if err := m.FetchByInstanceID(d); err != nil {
		return nil, status.CheckDatabaseError(err, "FetchInstanceByInstanceID")
	}

	state, ok := vm.GetInstanceState(instanceID.String())
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
	fmt.Println("---------------")

	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.Instance{}

	list, err := m.List(d, nil)
	if err != nil {
		return err
	}
	for _, i := range list {
		if i.State == enums.INSTANCE_STATE__CREATED || i.State == enums.INSTANCE_STATE__STARTED {
			err = vm.NewInstanceWithID(i.Path, i.InstanceID.String(), common.DefaultInstanceOptionSetter)
			if err != nil {
				if err := i.DeleteByInstanceID(d); err != nil {
					return err
				}
				if err := (&models.Applet{RelApplet: models.RelApplet{AppletID: i.AppletID}}).
					DeleteByAppletID(d); err != nil {
					return err
				}
				l.WithValues("instance", i.InstanceID, "applet", i.AppletID).
					Warn(errors.New("start failed and removed"))
				return nil
			}
			m.State = enums.INSTANCE_STATE__CREATED
		}
		if i.State == enums.INSTANCE_STATE__STARTED {
			_ = ControlInstance(ctx, i.InstanceID, enums.DEPLOY_CMD__START)
		}
	}
	return nil
}
