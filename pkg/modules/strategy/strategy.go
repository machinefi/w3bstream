package strategy

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

// TODO
func FilterStrategy(ctx context.Context, channel, eventType string) {}

type InstanceHandler struct {
	AppletID   types.SFID
	InstanceID types.SFID
	Handler    string
}

func FindStrategyInstances(ctx context.Context, prjName string, eventType string) ([]*InstanceHandler, error) {
	l := types.MustLoggerFromContext(ctx)
	d := types.MustMgrDBExecutorFromContext(ctx)

	_, l = l.Start(ctx, "FindStrategyInstances")
	defer l.End()

	l = l.WithValues("project", prjName, "event_type", eventType)

	mProject := &models.Project{ProjectName: models.ProjectName{Name: prjName}}

	if err := mProject.FetchByName(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "FetchProjectByName")
	}

	mStrategy := &models.Strategy{}

	strategies, err := mStrategy.List(d,
		builder.And(
			mStrategy.ColProjectID().Eq(mProject.ProjectID),
			builder.Or(
				mStrategy.ColEventType().Eq(eventType),
				mStrategy.ColEventType().Eq(enums.EVENTTYPEDEFAULT),
			),
		),
	)
	if err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "ListStrategy")
	}

	if len(strategies) == 0 {
		l.Warn(errors.New("strategy not found"))
		return nil, status.NotFound.StatusErr().WithDesc("not found strategy")
	}
	strategiesMap := make(map[types.SFID]*models.Strategy)
	for i := range strategies {
		strategiesMap[strategies[i].AppletID] = &strategies[i]
	}

	appletIDs := make(types.SFIDs, 0, len(strategies))

	for i := range strategies {
		appletIDs = append(appletIDs, strategies[i].AppletID)
	}

	mInstance := &models.Instance{}

	instances, err := mInstance.List(d, mInstance.ColAppletID().In(appletIDs))
	if err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "ListInstances")
	}

	handlers := make([]*InstanceHandler, 0)

	for _, instance := range instances {
		handlers = append(handlers, &InstanceHandler{
			AppletID:   instance.AppletID,
			InstanceID: instance.InstanceID,
			Handler:    strategiesMap[instance.AppletID].Handler,
		})
	}
	return handlers, nil
}

func Update(ctx context.Context, r *UpdateReq) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	sty := types.MustStrategyFromContext(ctx)

	sty.RelApplet = r.RelApplet
	sty.StrategyInfo = r.StrategyInfo

	if err := sty.UpdateByStrategyID(d); err != nil {
		if sqlx.DBErr(err).IsConflict() {
			return status.StrategyIsExists.StatusErr().WithDesc(
				fmt.Sprintf(
					"[prj: %s] [app: %s] [type: %s] [hdl: %s]",
					sty.ProjectID, sty.AppletID, sty.EventType, sty.Handler,
				),
			)
		}
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func GetBySFID(ctx context.Context, id types.SFID) (*models.Strategy, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Strategy{RelStrategy: models.RelStrategy{StrategyID: id}}

	if err := m.FetchByStrategyID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.StrategyNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func List(ctx context.Context, r *ListReq) (*ListRsp, error) {
	var (
		d    = types.MustMgrDBExecutorFromContext(ctx)
		prj  = types.MustProjectFromContext(ctx)
		sty  = &models.Strategy{}
		err  error
		ret  = &ListRsp{}
		cond = r.Condition(prj.ProjectID)
		adds = r.Additions()
	)

	ret.Data, err = sty.List(d, cond, adds...)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	ret.Total, err = sty.Count(d, cond)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return ret, nil
}

func ListDetail(ctx context.Context, r *ListReq) (*ListDetailRsp, error) {
	var (
		d   = types.MustMgrDBExecutorFromContext(ctx)
		sty = &models.Strategy{}
		app = &models.Applet{}
		ins = &models.Instance{}
		prj = types.MustProjectFromContext(ctx)
		ret = &ListDetailRsp{}

		cond = r.Condition(prj.ProjectID)
		adds = r.Additions()
	)

	expr := builder.Select(builder.MultiWith(",",
		builder.Alias(sty.ColStrategyID(), "f_sty_id"),
		builder.Alias(sty.ColProjectID(), "f_prj_id"),
		builder.Alias(prj.ColName(), "f_prj_name"),
		builder.Alias(sty.ColAppletID(), "f_app_id"),
		builder.Alias(app.ColName(), "f_app_name"),
		builder.Alias(ins.ColInstanceID(), "f_ins_id"),
		builder.Alias(sty.ColEventType(), "f_event_type"),
		builder.Alias(sty.ColHandler(), "f_handler"),
		builder.Alias(sty.ColUpdatedAt(), "f_updated_at"),
		builder.Alias(sty.ColCreatedAt(), "f_created_at"),
	)).From(
		d.T(sty),
		append([]builder.Addition{
			builder.LeftJoin(d.T(app)).On(sty.ColAppletID().Eq(app.ColAppletID())),
			builder.LeftJoin(d.T(prj)).On(sty.ColProjectID().Eq(prj.ColProjectID())),
			builder.Where(cond),
		}, adds...)...,
	)
	err := d.QueryAndScan(expr, &ret.Data)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	if ret.Total, err = sty.Count(d, cond); err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return ret, nil
}

func Remove(ctx context.Context) error {
	m := types.MustStrategyFromContext(ctx)

	if err := m.DeleteByStrategyID(types.MustMgrDBExecutorFromContext(ctx)); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func BatchRemove(ctx context.Context, r *DataListParam) error {
	var (
		d   = types.MustMgrDBExecutorFromContext(ctx)
		m   = &models.Strategy{}
		prj = types.MustProjectFromContext(ctx)
	)

	_, err := d.Exec(builder.Delete().From(
		d.T(m),
		builder.Where(r.Condition(prj.ProjectID)),
	))
	if err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func BatchCreate(ctx context.Context, sty []models.Strategy) error {
	if len(sty) == 0 {
		return nil
	}

	return sqlx.NewTasks(types.MustMgrDBExecutorFromContext(ctx)).With(
		func(d sqlx.DBExecutor) error {
			for i := range sty {
				s := &sty[i]
				err := s.Create(d)
				if err != nil {
					if sqlx.DBErr(err).IsConflict() {
						return status.StrategyIsExists.StatusErr().WithDesc(
							fmt.Sprintf(
								"[prj: %s] [app: %s] [type: %s] [hdl: %s]",
								s.ProjectID, s.AppletID, s.EventType, s.Handler,
							),
						)
					}
					return status.DatabaseError.StatusErr().WithDesc(err.Error())
				}
			}
			return nil
		},
	).Do()
}
