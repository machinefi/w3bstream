package deploy

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
	"github.com/pkg/errors"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateInstanceReq struct {
	Cache *wasm.Cache `json:"cache,omitempty"`
}

type CreateInstanceRsp struct {
	InstanceID    types.SFID          `json:"instanceID"`
	InstanceState enums.InstanceState `json:"instanceState"`
}

func CreateInstance(ctx context.Context, r *CreateInstanceReq) (*CreateInstanceRsp, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	app := types.MustAppletFromContext(ctx)
	res := types.MustResourceFromContext(ctx)
	ins := &models.Instance{
		RelInstance:  models.RelInstance{InstanceID: idg.MustGenSFID()},
		RelApplet:    models.RelApplet{AppletID: app.AppletID},
		InstanceInfo: models.InstanceInfo{State: enums.INSTANCE_STATE__CREATED},
	}

	_, l = l.Start(ctx, "CreateInstance")
	defer l.End()

	_ctx := context.Background()
	err := sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			if err := ins.Create(db); err != nil {
				return err
			}
			ctx = types.WithInstance(ctx, ins)
			return nil
		},
		func(db sqlx.DBExecutor) error {
			if r.Cache == nil {
				r.Cache = wasm.DefaultCache()
			}
			return config.CreateConfig(ctx, ins.InstanceID, r.Cache)
		},
		func(db sqlx.DBExecutor) error {
			var _err error
			_ctx, _err = WithInstanceRuntimeContext(ctx)
			return _err
		},
		func(db sqlx.DBExecutor) error {
			return vm.NewInstance(_ctx, res.Path, ins.InstanceID)
		},
	).Do()

	if err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err)
	}

	l.WithValues("instance", ins.InstanceID).Info("created")
	return &CreateInstanceRsp{
		InstanceID:    ins.InstanceID,
		InstanceState: ins.State,
	}, nil
}

func ControlInstance(ctx context.Context, instanceID types.SFID, cmd enums.DeployCmd) (err error) {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		l = types.MustLoggerFromContext(ctx)
		m = &models.Instance{RelInstance: models.RelInstance{InstanceID: instanceID}}
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

	if err = m.FetchByInstanceID(d); err != nil {
		l.Error(err)
		return status.CheckDatabaseError(err, "FetchByInstanceID")
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
		m.State = enums.INSTANCE_STATE_UNKNOWN
		if err = m.UpdateByInstanceID(d); err != nil {
			l.Error(err)
			return status.CheckDatabaseError(err, "UpdateInstanceByInstanceID")
		}
		return status.BadRequest.StatusErr().WithDesc("unknown deploy command")
	}
}

func GetInstanceByInstanceID(ctx context.Context, instanceID types.SFID) (*models.Instance, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
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
	d := types.MustMgrDBExecutorFromContext(ctx)
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
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.Instance{}

	_, l = l.Start(ctx, "StartInstances")
	defer l.End()

	list, err := m.List(d, nil)
	if err != nil {
		l.Error(err)
		return err
	}
	for i := range list {
		ins := &list[i]
		cmd := enums.DEPLOY_CMD_UNKNOWN
		l = l.WithValues(
			"instance", ins.InstanceID,
			"applet", ins.AppletID,
			"status", ins.State,
		)

		_ctx, err := WithInstanceRuntimeContext(types.WithInstance(ctx, ins))
		if err != nil {
			l.Error(err)
			continue
		}
		res := types.MustResourceFromContext(_ctx)
		if err = vm.NewInstance(_ctx, res.Path, ins.InstanceID); err != nil {
			l.Error(err)
			ins.State = enums.INSTANCE_STATE_UNKNOWN
		}
		switch ins.State {
		case enums.INSTANCE_STATE__CREATED:
			continue
		case enums.INSTANCE_STATE__STARTED:
			cmd = enums.DEPLOY_CMD__START
		case enums.INSTANCE_STATE__STOPPED:
			cmd = enums.DEPLOY_CMD__STOP
		default:
			cmd = enums.DEPLOY_CMD_UNKNOWN
		}

		l = l.WithValues("cmd", cmd)
		if err = ControlInstance(ctx, ins.InstanceID, cmd); err != nil {
			l.Error(err)
		}
	}
	return nil
}
