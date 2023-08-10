package strategy

import (
	"context"
	"fmt"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func Update(ctx context.Context, id types.SFID, r *UpdateReq) (err error) {
	var m *models.Strategy

	return sqlx.NewTasks(types.MustMgrDBExecutorFromContext(ctx)).With(
		func(d sqlx.DBExecutor) error {
			ctx := types.WithMgrDBExecutor(ctx, d)
			m, _ = types.StrategyFromContext(ctx)
			if m == nil || m.StrategyID != id {
				m, err = GetBySFID(ctx, id)
			}
			return err
		},
		func(d sqlx.DBExecutor) error {
			m.RelApplet = r.RelApplet
			m.StrategyInfo = r.StrategyInfo
			if err = m.UpdateByStrategyID(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.StrategyConflict.StatusErr().WithDesc(
						fmt.Sprintf(
							"[prj: %s] [app: %s] [type: %s] [hdl: %s]",
							m.ProjectID, m.AppletID, m.EventType, m.Handler,
						),
					)
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
	).Do()
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
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.Strategy{}

		err  error
		ret  = &ListRsp{}
		cond = r.Condition()
		adds = r.Additions()
	)

	ret.Data, err = m.List(d, cond, adds...)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	ret.Total, err = m.Count(d, cond)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return ret, nil
}

func ListByCond(ctx context.Context, r *CondArgs, adds ...builder.Addition) ([]models.Strategy, error) {
	data, err := (&models.Strategy{}).List(
		types.MustMgrDBExecutorFromContext(ctx),
		r.Condition(),
		adds...,
	)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return data, nil
}

func ListDetailByCond(ctx context.Context, r *CondArgs, adds ...builder.Addition) (data []*Detail, err error) {
	var (
		d    = types.MustMgrDBExecutorFromContext(ctx)
		sty  = &models.Strategy{}
		app  = &models.Applet{}
		ins  = &models.Instance{}
		prj  = &models.Project{}
		cond = r.Condition()
	)

	expr := builder.Select(builder.MultiWith(",",
		builder.Alias(prj.ColName(), "f_prj_name"),
		builder.Alias(sty.ColAppletID(), "f_app_id"),
		builder.Alias(app.ColName(), "f_app_name"),
		builder.Alias(ins.ColInstanceID(), "f_ins_id"),
		builder.Alias(sty.ColHandler(), "f_hdl"),
		builder.Alias(sty.ColEventType(), "f_evt"),
		builder.Alias(sty.ColAutoCollectMetric(), "f_auto_collect"),
		builder.Alias(sty.ColUpdatedAt(), "f_updated_at"),
		builder.Alias(sty.ColCreatedAt(), "f_created_at"),
	)).From(
		d.T(sty),
		append([]builder.Addition{
			builder.LeftJoin(d.T(app)).On(sty.ColAppletID().Eq(app.ColAppletID())),
			builder.LeftJoin(d.T(prj)).On(sty.ColProjectID().Eq(prj.ColProjectID())),
			builder.LeftJoin(d.T(ins)).On(sty.ColAppletID().Eq(ins.ColAppletID())),
			builder.Where(cond),
		}, adds...)...,
	)
	err = d.QueryAndScan(expr, &data)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return
}

func ListDetail(ctx context.Context, r *ListReq) (*ListDetailRsp, error) {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.Strategy{}

		err error
		ret = &ListDetailRsp{}
	)
	ret.Data, err = ListDetailByCond(ctx, &r.CondArgs, r.Addition())
	if ret.Total, err = m.Count(d, r.Condition()); err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return ret, nil
}

func RemoveBySFID(ctx context.Context, id types.SFID) error {
	m := &models.Strategy{RelStrategy: models.RelStrategy{StrategyID: id}}

	if err := m.DeleteByStrategyID(types.MustMgrDBExecutorFromContext(ctx)); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func Remove(ctx context.Context, r *CondArgs) error {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.Strategy{}
	)

	_, err := d.Exec(builder.Delete().From(
		d.T(m),
		builder.Where(r.Condition()),
	))
	if err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func Create(ctx context.Context, r *CreateReq) (*models.Strategy, error) {
	var (
		idg = confid.MustNewSFIDGenerator()
		app *models.Applet
		sty = &models.Strategy{
			RelApplet:    r.RelApplet,
			StrategyInfo: r.StrategyInfo,
		}
	)

	err := sqlx.NewTasks(types.MustMgrDBExecutorFromContext(ctx)).With(
		func(d sqlx.DBExecutor) error {
			app, _ = types.AppletFromContext(ctx)
			if app == nil || app.AppletID != sty.AppletID {
				app = &models.Applet{
					RelApplet: models.RelApplet{AppletID: sty.AppletID},
				}
				if err := app.FetchByAppletID(d); err != nil {
					if sqlx.DBErr(err).IsNotFound() {
						return status.AppletNotFound
					}
					return status.DatabaseError.StatusErr().WithDesc(err.Error())
				}
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			sty.ProjectID = app.ProjectID
			sty.StrategyID = idg.MustGenSFID()
			if err := sty.Create(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.StrategyConflict.StatusErr().WithDesc(
						fmt.Sprintf(
							"[prj: %s] [app: %s] [type: %s] [hdl: %s]",
							sty.ProjectID, sty.AppletID, sty.EventType, sty.Handler,
						),
					)
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
	).Do()
	if err != nil {
		return nil, err
	}
	return sty, nil
}

func BatchCreate(ctx context.Context, sty []models.Strategy) error {
	if len(sty) == 0 {
		return nil
	}

	return sqlx.NewTasks(types.MustMgrDBExecutorFromContext(ctx)).With(
		func(d sqlx.DBExecutor) error {
			for i := range sty {
				s := &sty[i]
				if err := s.Create(d); err != nil {
					if sqlx.DBErr(err).IsConflict() {
						return status.StrategyConflict.StatusErr().WithDesc(
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

func FilterByProjectAndEvent(ctx context.Context, id types.SFID, tpe string) ([]*types.StrategyResult, error) {
	data, err := ListDetailByCond(ctx, &CondArgs{
		ProjectID: id, EventTypes: []string{tpe}},
	)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		data, err = ListDetailByCond(ctx, &CondArgs{
			ProjectID: id, EventTypes: []string{enums.EVENTTYPEDEFAULT}},
		)
		if err != nil {
			return nil, err
		}
	}

	results := make([]*types.StrategyResult, 0, len(data))
	for i := range data {
		results = append(results, &data[i].StrategyResult)
	}

	return results, nil
}
