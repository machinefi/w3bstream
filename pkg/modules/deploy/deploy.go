package deploy

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/machinefi/w3bstream/pkg/types"
)

func Init(ctx context.Context) error {
	var (
		d   = types.MustMgrDBExecutorFromContext(ctx)
		ins = &models.Instance{}
		app *models.Applet
		res *models.Resource
	)

	_, l := types.MustLoggerFromContext(ctx).Start(ctx, "deploy.Init")
	defer l.End()

	list, err := ins.List(d, nil)
	if err != nil {
		return err
	}
	for i := range list {
		ins = &list[i]
		l = l.WithValues("ins", ins.InstanceID, "app", ins.AppletID)

		app = &models.Applet{RelApplet: models.RelApplet{AppletID: ins.AppletID}}
		err = app.FetchByAppletID(d)
		if err != nil {
			l.Warn(err)
			continue
		}

		l = l.WithValues("res", app.ResourceID)
		res = &models.Resource{RelResource: models.RelResource{ResourceID: app.ResourceID}}
		err = app.FetchByAppletID(d)
		if err != nil {
			l.Warn(err)
			continue
		}

		ctx = contextx.WithContextCompose(
			types.WithResourceContext(res),
			types.WithAppletContext(app),
		)(ctx)

		if state := ins.State; state == enums.INSTANCE_STATE__STARTED ||
			state == enums.INSTANCE_STATE__STOPPED {
			ins, err = Upsert(ctx, state, ins.InstanceID)
			if err != nil {
				l.Warn(err)
			} else {
				l.WithValues("state", ins.State)
			}
		}
	}
	return nil
}

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

func ListWithCond(ctx context.Context, r *CondArgs) (data []models.Instance, err error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Instance{}

	if r.ProjectID == 0 {
		data, err = m.List(d, r.Condition())
	} else {
		app := &models.Applet{}
		err = d.QueryAndScan(
			builder.Select(nil).From(
				d.T(m),
				builder.LeftJoin(d.T(app)).On(m.ColAppletID().Eq(app.ColAppletID())),
				builder.Where(r.Condition()),
			), &data,
		)
	}
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return data, nil
}

func RemoveBySFID(ctx context.Context, id types.SFID) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Instance{RelInstance: models.RelInstance{InstanceID: id}}

	return sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			if err := m.DeleteByInstanceID(d); err != nil {
				return status.DatabaseError.StatusErr().
					WithDesc(errors.Wrap(err, id.String()).Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			return config.Remove(ctx, &config.CondArgs{RelIDs: []types.SFID{id}})
		},
		func(d sqlx.DBExecutor) error {
			if err := vm.DelInstance(ctx, m.InstanceID); err != nil {
				// Warn
			}
			return nil
		},
	).Do()
}

func RemoveByAppletSFID(ctx context.Context, id types.SFID) (err error) {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m *models.Instance
	)

	return sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			m, err = GetByAppletSFID(ctx, id)
			return err
		},
		func(d sqlx.DBExecutor) error {
			return RemoveBySFID(ctx, m.InstanceID)
		},
	).Do()
}

func Remove(ctx context.Context, r *CondArgs) error {
	var (
		lst []models.Instance
		err error
	)

	return sqlx.NewTasks(types.MustMgrDBExecutorFromContext(ctx)).With(
		func(db sqlx.DBExecutor) error {
			lst, err = ListWithCond(ctx, r)
			return err
		},
		func(db sqlx.DBExecutor) error {
			for i := range lst {
				err = RemoveBySFID(ctx, lst[i].InstanceID)
				if err != nil {
					return err
				}
			}
			return nil
		},
	).Do()
}

func UpsertByCode(ctx context.Context, code []byte, state enums.InstanceState, old ...types.SFID) (*models.Instance, error) {
	var (
		id        types.SFID
		forUpdate = false
	)

	if len(old) > 0 && old[0] != 0 {
		forUpdate = true
		id = old[0]
	} else {
		id = confid.MustSFIDGeneratorFromContext(ctx).MustGenSFID()
	}
	app := types.MustAppletFromContext(ctx)
	ins := &models.Instance{
		RelInstance:  models.RelInstance{InstanceID: id},
		RelApplet:    models.RelApplet{AppletID: app.AppletID},
		InstanceInfo: models.InstanceInfo{State: state},
	}

	err := sqlx.NewTasks(types.MustMgrDBExecutorFromContext(ctx)).With(
		func(d sqlx.DBExecutor) error {
			if !forUpdate {
				return nil
			}
			if err := ins.UpdateByInstanceID(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.MultiInstanceDeployed.StatusErr().
						WithDesc(app.AppletID.String())
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			if forUpdate {
				return nil
			}
			if err := ins.Create(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.MultiInstanceDeployed.StatusErr().
						WithDesc(app.AppletID.String())
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			if forUpdate {
				if err := vm.DelInstance(ctx, ins.InstanceID); err != nil {
					// Warn
				}
			}
			_ctx, err := WithInstanceRuntimeContext(types.WithInstance(ctx, ins))
			if err != nil {
				return err
			}
			// TODO should below actions be in a critical section?
			if err = vm.NewInstance(_ctx, code, id, state); err != nil {
				return status.CreateInstanceFailed.StatusErr().WithDesc(err.Error())
			}
			ins.State, _ = vm.GetInstanceState(ins.InstanceID)
			if ins.State != state {
				// Warn
			}
			return nil
		},
	).Do()
	if err != nil {
		return nil, err
	}
	return ins, nil
}

func Upsert(ctx context.Context, state enums.InstanceState, old ...types.SFID) (*models.Instance, error) {
	res := types.MustResourceFromContext(ctx)

	code, err := resource.GetContentBySFID(ctx, res.ResourceID)
	if err != nil {
		return nil, err
	}

	return UpsertByCode(ctx, code, state, old...)
}

func Create(ctx context.Context, r *CreateReq) (*models.Instance, error) {
	var (
		idg = confid.MustSFIDGeneratorFromContext(ctx)
		app = types.MustAppletFromContext(ctx)
		ins *models.Instance
		err error
	)

	err = sqlx.NewTasks(types.MustMgrDBExecutorFromContext(ctx)).With(
		func(d sqlx.DBExecutor) error {
			ins = &models.Instance{
				RelInstance:  models.RelInstance{InstanceID: idg.MustGenSFID()},
				RelApplet:    models.RelApplet{AppletID: app.AppletID},
				InstanceInfo: models.InstanceInfo{State: enums.INSTANCE_STATE__CREATED},
			}
			if err = ins.Create(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.MultiInstanceDeployed.StatusErr().
						WithDesc(app.AppletID.String())
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			_, err = config.Create(ctx, ins.InstanceID, r.Cache)
			return nil
		},
	).Do()
	if err != nil {
		return nil, err
	}
	return ins, nil
}

func Deploy(ctx context.Context, cmd enums.DeployCmd) (err error) {
	var m = types.MustInstanceFromContext(ctx)

	switch cmd {
	case enums.DEPLOY_CMD__STOP:
		m.State = enums.INSTANCE_STATE__STOPPED
	case enums.DEPLOY_CMD__START:
		m.State = enums.INSTANCE_STATE__STARTED
	case enums.DEPLOY_CMD__REMOVE:
		m.State = enums.INSTANCE_STATE__CREATED
	default:
		return status.UnknownDeployCommand.StatusErr().
			WithDesc(strconv.Itoa(int(cmd)))
	}

	return sqlx.NewTasks(types.MustMgrDBExecutorFromContext(ctx)).With(
		func(d sqlx.DBExecutor) error {
			if err = m.UpdateByInstanceID(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.MultiInstanceDeployed.StatusErr().
						WithDesc(m.AppletID.String())
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			switch m.State {
			case enums.INSTANCE_STATE__STOPPED:
				err = vm.StopInstance(ctx, m.InstanceID)
			case enums.INSTANCE_STATE__STARTED:
				err = vm.StartInstance(ctx, m.InstanceID)
			case enums.INSTANCE_STATE__CREATED:
				err = vm.DelInstance(ctx, m.InstanceID)
			}
			if err != nil {
				// Warn
			}
			return nil
		},
	).Do()
}
