package deploy

import (
	"context"
	"fmt"
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
	"github.com/machinefi/w3bstream/pkg/modules/job"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/machinefi/w3bstream/pkg/modules/wasmlog"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm/kvdb"
)

func Init(ctx context.Context) error {
	var (
		d   = types.MustMgrDBExecutorFromContext(ctx)
		rDB = kvdb.MustRedisDBKeyFromContext(ctx)

		ins     = &models.Instance{}
		traffic = &models.TrafficLimit{}
		prj     *models.Project
		app     *models.Applet
		res     *models.Resource
	)

	_, l := types.MustLoggerFromContext(ctx).Start(ctx, "deploy.Init")
	defer l.End()

	trafficList, err := traffic.List(d, nil)
	if err != nil {
		l.Error(err)
		return err
	}
	for i := range trafficList {
		traffic = &trafficList[i]
		prj = &models.Project{RelProject: models.RelProject{ProjectID: traffic.ProjectID}}
		err = prj.FetchByProjectID(d)
		if err != nil {
			l.Warn(err)
			continue
		}
		projectKey := fmt.Sprintf("%s::%s", prj.Name, traffic.ApiType.String())
		valByte, err := rDB.GetKey(projectKey)
		if err != nil {
			l.Warn(err)
			continue
		}
		if valByte == nil {
			err = rDB.SetKeyWithEX(projectKey,
				[]byte(strconv.Itoa(traffic.Threshold)), 31622400)
		}
		t := job.NewTrafficTaskWithPrjKey(projectKey, *traffic)
		job.Dispatch(ctx, t)
	}

	list, err := ins.List(d, nil)
	if err != nil {
		l.Error(err)
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
			ins, err = Upsert(ctx, nil, state, ins.InstanceID)
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
			ctx := types.WithMgrDBExecutor(ctx, d)
			return config.Remove(ctx, &config.CondArgs{RelIDs: []types.SFID{id}})
		},
		func(d sqlx.DBExecutor) error {
			if err := vm.DelInstance(ctx, m.InstanceID); err != nil {
				// Warn
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			ctx := types.WithMgrDBExecutor(ctx, d)
			return wasmlog.Remove(ctx, &wasmlog.CondArgs{InstanceID: m.InstanceID})
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
			ctx := types.WithMgrDBExecutor(ctx, d)
			m, err = GetByAppletSFID(ctx, id)
			return err
		},
		func(d sqlx.DBExecutor) error {
			ctx := types.WithMgrDBExecutor(ctx, d)
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
			ctx := types.WithMgrDBExecutor(ctx, db)
			lst, err = ListWithCond(ctx, r)
			return err
		},
		func(db sqlx.DBExecutor) error {
			ctx := types.WithMgrDBExecutor(ctx, db)
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

// UpsertByCode upsert instance and its config, and deploy wasm if needed
func UpsertByCode(ctx context.Context, r *CreateReq, code []byte, state enums.InstanceState, old ...types.SFID) (*models.Instance, error) {
	var (
		idg       = confid.MustSFIDGeneratorFromContext(ctx)
		forUpdate = false
	)

	app := types.MustAppletFromContext(ctx)
	ins := &models.Instance{}

	err := sqlx.NewTasks(types.MustMgrDBExecutorFromContext(ctx)).With(
		func(d sqlx.DBExecutor) error {
			ins.AppletID = app.AppletID
			if err := ins.FetchByAppletID(d); err != nil {
				if sqlx.DBErr(err).IsNotFound() {
					forUpdate = false
					ins.InstanceID = idg.MustGenSFID()
					ins.State = state
					return nil
				} else {
					return status.DatabaseError.StatusErr().WithDesc(err.Error())
				}
			}
			if len(old) > 0 && old[0] != ins.InstanceID {
				return status.InvalidAppletContext.StatusErr().WithDesc(
					fmt.Sprintf("database: %v arg: %v", ins.InstanceID, old[0]),
				)
			}
			ins.State = state
			forUpdate = true
			return nil
		},
		func(d sqlx.DBExecutor) error {
			if forUpdate {
				if err := ins.UpdateByInstanceID(d); err != nil {
					if sqlx.DBErr(err).IsConflict() {
						return status.MultiInstanceDeployed.StatusErr().
							WithDesc(app.AppletID.String())
					}
					return status.DatabaseError.StatusErr().WithDesc(err.Error())
				}
			} else {
				if err := ins.Create(d); err != nil {
					if sqlx.DBErr(err).IsConflict() {
						return status.MultiInstanceDeployed.StatusErr().
							WithDesc(app.AppletID.String())
					}
					return status.DatabaseError.StatusErr().WithDesc(err.Error())
				}
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			ctx := types.WithMgrDBExecutor(ctx, db)
			if r != nil && r.Cache != nil {
				if err := config.Remove(ctx, &config.CondArgs{
					RelIDs: []types.SFID{ins.InstanceID},
				}); err != nil {
					return err
				}
				if _, err := config.Create(ctx, ins.InstanceID, r.Cache); err != nil {
					return err
				}
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
			if err = vm.NewInstance(_ctx, code, ins.InstanceID, state); err != nil {
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

func Upsert(ctx context.Context, r *CreateReq, state enums.InstanceState, old ...types.SFID) (*models.Instance, error) {
	res := types.MustResourceFromContext(ctx)

	_, code, err := resource.GetContentBySFID(ctx, res.ResourceID)
	if err != nil {
		return nil, err
	}

	return UpsertByCode(ctx, r, code, state, old...)
}

func Deploy(ctx context.Context, cmd enums.DeployCmd) (err error) {
	var m = types.MustInstanceFromContext(ctx)

	switch cmd {
	case enums.DEPLOY_CMD__HUNGUP:
		m.State = enums.INSTANCE_STATE__STOPPED
	case enums.DEPLOY_CMD__START:
		m.State = enums.INSTANCE_STATE__STARTED
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
			}
			if err != nil {
				// Warn
			}
			return nil
		},
	).Do()
}
