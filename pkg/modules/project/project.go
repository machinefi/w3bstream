// project management

package project

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/mq"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateProjectReq = models.ProjectInfo

func CreateProject(ctx context.Context, r *CreateProjectReq, hdl mq.OnMessage) (*models.Project, error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	a := middleware.CurrentAccountFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "CreateProject")
	defer l.End()

	m := &models.Project{
		RelProject:  models.RelProject{ProjectID: idg.MustGenSFID()},
		RelAccount:  models.RelAccount{AccountID: a.AccountID},
		ProjectInfo: *r,
	}

	if err := mq.CreateChannel(ctx, m.Name, hdl); err != nil {
		l.Error(err)
		return nil, status.InternalServerError.StatusErr().
			WithDesc(fmt.Sprintf("create channel: [project:%s] [err:%v]", m.Name, err))
	}

	if err := m.Create(d); err != nil {
		l.Error(err)
		return nil, err
	}

	return m, nil
}

type ListProjectReq struct {
	accountID  types.SFID
	IDs        []uint64     `in:"query" name:"ids,omitempty"`
	ProjectIDs []types.SFID `in:"query" name:"projectIDs,omitempty"`
	Names      []string     `in:"query" name:"names,omitempty"`
	NameLike   string       `in:"query" name:"name,omitempty"`
	datatypes.Pager
}

func (r *ListProjectReq) SetCurrentAccount(accountID types.SFID) {
	r.accountID = accountID
}

func (r *ListProjectReq) Condition() builder.SqlCondition {
	var (
		m  = &models.Project{}
		cs []builder.SqlCondition
	)

	cs = append(cs, m.ColAccountID().Eq(r.accountID))
	if len(r.IDs) > 0 {
		cs = append(cs, m.ColID().In(r.IDs))
	}
	if len(r.ProjectIDs) > 0 {
		cs = append(cs, m.ColProjectID().In(r.ProjectIDs))
	}
	if len(r.Names) > 0 {
		cs = append(cs, m.ColName().In(r.Names))
	}
	if r.NameLike != "" {
		cs = append(cs, m.ColName().Like(r.NameLike))
	}

	return builder.And(cs...)
}

func (r *ListProjectReq) Additions() builder.Additions {
	m := &models.Project{}
	return builder.Additions{
		builder.OrderBy(builder.DescOrder(m.ColCreatedAt())),
		r.Pager.Addition(),
	}
}

type ListProjectRsp struct {
	Data  []Detail `json:"data"`  // Data project data list
	Total int64    `json:"total"` // Total project count under current user
}

type Detail struct {
	ProjectID   types.SFID     `json:"projectID"`
	ProjectName string         `json:"projectName"`
	Applets     []AppletDetail `json:"applets,omitempty"`
	datatypes.OperationTimes
}

type AppletDetail struct {
	AppletID        types.SFID          `json:"appletID"`
	AppletName      string              `json:"appletName"`
	InstanceID      types.SFID          `json:"instanceID,omitempty"`
	InstanceState   enums.InstanceState `json:"instanceState,omitempty"`
	InstanceVMState enums.InstanceState `json:"instanceStateVM,omitempty"`
}

type detail struct {
	ProjectID     types.SFID          `db:"f_project_id"`
	ProjectName   string              `db:"f_project_name"`
	AppletID      types.SFID          `db:"f_applet_id"`
	AppletName    string              `db:"f_applet_name"`
	InstanceID    types.SFID          `db:"f_instance_id"`
	InstanceState enums.InstanceState `db:"f_instance_state"`
	datatypes.OperationTimes
}

func ListProject(ctx context.Context, r *ListProjectReq) (*ListProjectRsp, error) {
	var (
		d = types.MustDBExecutorFromContext(ctx)
		l = types.MustLoggerFromContext(ctx)

		ret  = &ListProjectRsp{}
		err  error
		cond = r.Condition()

		mProject  = &models.Project{}
		mApplet   = &models.Applet{}
		mInstance = &models.Instance{}
	)

	_, l = l.Start(ctx, "ListProject")
	defer l.End()

	ret.Total, err = mProject.Count(d, cond)
	if err != nil {
		return nil, status.CheckDatabaseError(err, "CountProject")
	}

	details := make([]detail, 0)

	err = d.QueryAndScan(
		builder.Select(
			builder.MultiWith(
				",",
				builder.Alias(mProject.ColProjectID(), "f_project_id"),
				builder.Alias(mProject.ColName(), "f_project_name"),
				builder.Alias(mApplet.ColAppletID(), "f_applet_id"),
				builder.Alias(mApplet.ColName(), "f_applet_name"),
				builder.Alias(mInstance.ColInstanceID(), "f_instance_id"),
				builder.Alias(mInstance.ColState(), "f_instance_state"),
				builder.Alias(mProject.ColCreatedAt(), "f_created_at"),
				builder.Alias(mProject.ColUpdatedAt(), "f_updated_at"),
			),
		).From(
			d.T(mProject),
			builder.LeftJoin(d.T(mApplet)).
				On(mProject.ColProjectID().Eq(mApplet.ColProjectID())),
			builder.LeftJoin(d.T(mInstance)).
				On(mApplet.ColAppletID().Eq(mInstance.ColAppletID())),
			builder.Where(cond),
			builder.OrderBy(
				builder.DescOrder(mProject.ColCreatedAt()),
				builder.AscOrder(mApplet.ColName()),
			),
			r.Pager.Addition(),
		),
		&details,
	)
	if err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "ListProject")
	}

	detailsMap := make(map[types.SFID][]*detail)
	for i := range details {
		prjID := details[i].ProjectID
		detailsMap[prjID] = append(detailsMap[prjID], &details[i])
	}

	for prjID, vmap := range detailsMap {
		appletDetails := make([]AppletDetail, 0, len(vmap))
		for _, v := range vmap {
			if v.AppletID == 0 {
				continue
			}
			state, ok := vm.GetInstanceState(v.InstanceID)
			if !ok {
				l.Warn(errors.New("instance not found in vms"))
			}
			appletDetails = append(appletDetails, AppletDetail{
				AppletID:        v.AppletID,
				AppletName:      v.AppletName,
				InstanceID:      v.InstanceID,
				InstanceState:   v.InstanceState,
				InstanceVMState: state,
			})
		}
		if len(appletDetails) == 0 {
			appletDetails = nil
		}
		ret.Data = append(ret.Data, Detail{
			ProjectID:   prjID,
			ProjectName: vmap[0].ProjectName,
			Applets:     appletDetails,
			OperationTimes: datatypes.OperationTimes{
				CreatedAt: vmap[0].CreatedAt,
				UpdatedAt: vmap[0].UpdatedAt,
			},
		})
	}

	return ret, nil
}

func GetProjectByProjectID(ctx context.Context, prjID types.SFID) (*Detail, error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	ca := middleware.CurrentAccountFromContext(ctx)

	_, l = l.Start(ctx, "GetProjectByProjectID")
	defer l.End()

	_, err := ca.ValidateProjectPerm(ctx, prjID)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	m := &models.Project{RelProject: models.RelProject{ProjectID: prjID}}

	if err = m.FetchByProjectID(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "GetProjectByProjectID")
	}

	ret, err := ListProject(ctx, &ListProjectReq{
		accountID:  ca.AccountID,
		ProjectIDs: []types.SFID{prjID},
	})

	if err != nil {
		l.Error(err)
		return nil, err
	}

	if len(ret.Data) == 0 {
		l.Warn(errors.New("project not found"))
		return nil, status.NotFound
	}

	return &ret.Data[0], nil
}

func GetProjectByProjectName(ctx context.Context, prjName string) (*models.Project, error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.Project{ProjectInfo: models.ProjectInfo{Name: prjName}}

	_, l = l.Start(ctx, "GetProjectByProjectName")
	defer l.End()

	if err := m.FetchByName(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "GetProjectByProjectName")
	}

	return m, nil
}

func DeleteProject(_ context.Context, _ string) error {
	// TODO
	// same as RemoveProjectByProjectID?  by zhiwei
	return nil
}

func InitChannels(ctx context.Context, hdl mq.OnMessage) error {
	l := types.MustLoggerFromContext(ctx)
	d := types.MustDBExecutorFromContext(ctx)
	m := &models.Project{}

	_, l = l.Start(ctx, "InitChannels")
	defer l.End()

	lst, err := m.List(d, nil)
	if err != nil {
		l.Error(err)
		return err
	}

	for i := range lst {
		v := &lst[i]
		err = mq.CreateChannel(ctx, v.Name, hdl)
		if err != nil {
			err = errors.Errorf("create channel: [project:%s] [err:%v]", v.Name, err)
			l.Error(err)
			return err
		}
		l.WithValues("project_name", v.Name).Info("mqtt subscribe started")
	}
	return nil
}

func RemoveProjectByProjectID(ctx context.Context, prjID types.SFID) error {
	var (
		d          = types.MustDBExecutorFromContext(ctx)
		l          = types.MustLoggerFromContext(ctx)
		mProject   = &models.Project{}
		mStrategy  = &models.Strategy{}
		strategies []models.Strategy
		mPublisher = &models.Publisher{}
		publishers []models.Publisher
		mApplet    = &models.Applet{}
		applets    []models.Applet
		mInstance  = &models.Instance{}
		instances  []models.Instance
		err        error
	)

	_, l = l.Start(ctx, "RemoveProjectByProjectID")
	defer l.End()

	return sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			mProject.ProjectID = prjID
			err = mProject.FetchByProjectID(d)
			if err != nil {
				l.Error(err)
				return status.CheckDatabaseError(err, "fetch by project id")
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			mStrategy.ProjectID = prjID
			strategies, err = mStrategy.List(d, mStrategy.ColProjectID().Eq(prjID))
			if err != nil {
				l.Error(err)
				return status.CheckDatabaseError(err, "ListStrategiesByProjectID")
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			mPublisher.ProjectID = prjID
			publishers, err = mPublisher.List(d, mPublisher.ColProjectID().Eq(prjID))
			if err != nil {
				l.Error(err)
				return status.CheckDatabaseError(err, "ListPublishersByProjectID")
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			mApplet.ProjectID = prjID
			applets, err = mApplet.List(d, mApplet.ColProjectID().Eq(prjID))
			if err != nil {
				l.Error(err)
				return status.CheckDatabaseError(err, "ListAppletsByProjectID")
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			for _, app := range applets {
				mInstance.AppletID = app.AppletID
				if tmp, e := mInstance.List(d, mInstance.ColAppletID().Eq(app.AppletID)); e != nil {
					l.Error(err)
					return status.CheckDatabaseError(err, "ListByAppletID")
				} else {
					instances = append(instances, tmp...)
				}
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
			for _, app := range applets {
				err = app.DeleteByAppletID(d)
				if err != nil {
					l.Error(err)
					return status.CheckDatabaseError(err, "DeleteAppletByAppletID")
				}
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			for _, strategy := range strategies {
				err = strategy.DeleteByStrategyID(d)
				if err != nil {
					l.Error(err)
					return status.CheckDatabaseError(err, "DeleteStrategyByStrategyID")
				}
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			for _, publisher := range publishers {
				err = publisher.DeleteByPublisherID(d)
				if err != nil {
					l.Error(err)
					return status.CheckDatabaseError(err, "DeletePublisherByStrategyID")
				}
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			err = mProject.DeleteByProjectID(d)
			if err != nil {
				l.Error(err)
				return status.CheckDatabaseError(err, "DeleteProjectByProjectID")
			}
			return nil
		},
	).Do()
}
