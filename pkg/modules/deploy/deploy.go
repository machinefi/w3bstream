package deploy

import (
	"context"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/machinefi/w3bstream/pkg/types"
)

func GetBySFID(ctx context.Context, id types.SFID) (*models.Instance, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Instance{RelInstance: models.RelInstance{InstanceID: id}}

	if err := m.FetchByInstanceID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.InstanceNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	m.State, _ = vm.GetInstanceState(m.InstanceID)
	return m, nil
}

func GetByAppletSFID(ctx context.Context, id types.SFID) (*models.Instance, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Instance{RelApplet: models.RelApplet{AppletID: id}}

	if err := m.FetchByAppletID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.InstanceNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	m.State, _ = vm.GetInstanceState(m.InstanceID)
	return m, nil
}

// Create creates deploy instance
// required: Instance or Applet context
func Create(ctx context.Context, r *CreateReq) (*models.Instance, error) {
	var (
		d      = types.MustMgrDBExecutorFromContext(ctx)
		idg    = confid.MustSFIDGeneratorFromContext(ctx)
		app    = types.MustAppletFromContext(ctx)
		err    error
		ins, _ = types.InstanceFromContext(ctx)
		_ins   *models.Instance
	)

	if ins == nil {
		_ins = &models.Instance{
			RelInstance: models.RelInstance{InstanceID: idg.MustGenSFID()},
			RelApplet:   models.RelApplet{AppletID: app.AppletID},
		}
	} else {
		_ins = ins
	}

	err = sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			if ins != nil {
				return nil
			}
			if err = _ins.Create(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.MultiInstanceRunning
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			ctx = types.WithInstance(ctx, ins)
			return nil
		},
		func(d sqlx.DBExecutor) error {
			if r.CacheConfig != nil {
				return config.Create(ctx, _ins.InstanceID, r.CacheConfig)
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			state, err := Deploy(ctx, enums.DEPLOY_CMD__START)
			ins.State = state
			return err
		},
	).Do()

	if err != nil {
		return nil, err
	}
	return _ins, nil
}

func deploy(ctx context.Context, state enums.InstanceState, cmd enums.DeployCmd, id types.SFID) (enums.InstanceState, error) {
	var err error

	switch cmd {
	case enums.DEPLOY_CMD__START, enums.DEPLOY_CMD__RESTART:
		if state, _ = vm.GetInstanceState(id); state == enums.INSTANCE_STATE__STARTED {
			return state, status.StartInstanceFailed.StatusErr().WithDesc(err.Error())
		}
		_ = vm.DelInstance(ctx, id)
		var (
			code   []byte
			wsmctx context.Context
		)
		wsmctx, err = WithInstanceRuntimeContext(ctx)
		if err != nil {
			return state, err
		}
		if code, err = resource.Load(wsmctx); err != nil {
			return state, err
		}
		if err = vm.NewInstanceByCode(wsmctx, code, id); err != nil {
			return state, status.StartInstanceFailed.StatusErr().WithDesc(err.Error())
		}
		if err = vm.StartInstance(ctx, id); err != nil {
			return state, status.StartInstanceFailed.StatusErr().WithDesc(err.Error())
		}
		state, _ = vm.GetInstanceState(id)
		return state, nil
	case enums.DEPLOY_CMD__STOP:
		if err = vm.StopInstance(ctx, id); err != nil {
			return state, status.StopInstanceFailed.StatusErr().WithDesc(err.Error())
		}
		return enums.INSTANCE_STATE__STOPPED, nil
	}
	// Warn
	return state, status.InvalidDeployCommand
}

func Deploy(ctx context.Context, cmd enums.DeployCmd) (enums.InstanceState, error) {
	var (
		d     = types.MustMgrDBExecutorFromContext(ctx)
		ins   = types.MustInstanceFromContext(ctx)
		err   error
		state enums.InstanceState
	)

	err = sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			if ins.State == state {
				return nil
			}
			ins.State = state
			if err = ins.UpdateByInstanceID(d); err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			state, err = deploy(ctx, ins.State, cmd, ins.InstanceID)
			return err
		},
	).Do()
	return state, err
}

// TODO should impl Restart(), for refreshing wasm context

func Remove(ctx context.Context, keepConfig bool) error {
	l := types.MustLoggerFromContext(ctx)
	ins, _ := types.InstanceFromContext(ctx)
	if ins == nil {
		return nil
	}

	return sqlx.NewTasks(types.MustMgrDBExecutorFromContext(ctx)).With(
		func(d sqlx.DBExecutor) error {
			if err := ins.DeleteByInstanceID(d); err != nil {
				if sqlx.DBErr(err).IsNotFound() {
					return status.InstanceNotFound
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			if !keepConfig {
				return config.BatchRemove(
					ctx,
					&config.DataListParam{RelIDs: types.SFIDs{ins.InstanceID}},
				)
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			// TODO if DeleteInstance return nil when InstanceID is nonexistent.
			if err := vm.DelInstance(ctx, ins.InstanceID); err != nil {
				l.WithValues("app", ins.AppletID, "ins", ins.InstanceID).
					Warn(err)
				// return status.DeleteInstanceFailed.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
	).Do()
}

func Init(ctx context.Context) {
	_, l := conflog.FromContext(ctx).Start(ctx, "deploy.Init")
	defer l.End()

	data, err := (&models.Instance{}).List(types.MustMgrDBExecutorFromContext(ctx), nil)
	if err != nil {
		l.Panic(err)
	}
	for i := range data {
		ins := &data[i]
		l = l.WithValues("ins", ins.InstanceID, "app", ins.AppletID)

		if ins.State != enums.INSTANCE_STATE__STOPPED {
			state, err := Deploy(types.WithInstance(ctx, ins), enums.DEPLOY_CMD__START)
			if err != nil {
				l.Error(err)
				continue
			}
			l.WithValues("state", state).Info("started")
		}
	}
}
