package strategy

import (
	"context"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func ListResultsByProjectAndEventType(ctx context.Context, id types.SFID, et string) (ret []*types.StrategyResult, err error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	sty := &models.Strategy{}
	prj := &models.Project{}
	app := &models.Applet{}
	ins := &models.Instance{}

	exp := builder.Select(builder.Multi(
		builder.Alias(prj.ColAccountID(), "f_acc_id"),
		builder.Alias(sty.ColProjectID(), "f_prj_id"),
		builder.Alias(prj.ColName(), "f_prj_name"),
		builder.Alias(app.ColAppletID(), "f_app_id"),
		builder.Alias(app.ColName(), "f_app_name"),
		builder.Alias(ins.ColInstanceID(), "f_ins_id"),
		builder.Alias(sty.ColHandler(), "f_hdl"),
		builder.Alias(sty.ColEventType(), "f_evt"),
	)).From(
		d.T(sty),
		builder.LeftJoin(d.T(prj)).On(sty.ColProjectID().Eq(prj.ColProjectID())),
		builder.LeftJoin(d.T(app)).On(app.ColProjectID().Eq(prj.ColProjectID())),
		builder.LeftJoin(d.T(ins)).On(ins.ColAppletID().Eq(app.ColAppletID())),
		builder.Where(builder.And(
			sty.ColDeletedAt().Eq(0), prj.ColDeletedAt().Eq(0),
			ins.ColDeletedAt().Eq(0), app.ColDeletedAt().Eq(0),
			builder.Or(sty.ColEventType().Eq(et), sty.ColEventType().Eq(enums.EVENTTYPEDEFAULT)),
			ins.ColState().Eq(enums.INSTANCE_STATE__STARTED),
			prj.ColProjectID().Eq(id),
		)),
	)

	ret = make([]*types.StrategyResult, 0)
	if err = d.QueryAndScan(exp, &ret); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.StrategyNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return
}

type CreateStrategyBatchReq struct {
	Strategies []CreateStrategyReq `json:"strategies"`
}

type CreateStrategyReq struct {
	models.RelApplet
	models.StrategyInfo
}

func CreateStrategy(ctx context.Context, projectID types.SFID, r *CreateStrategyBatchReq) (err error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "CreateStrategy")
	defer l.End()

	err = sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			for i := range r.Strategies {
				if err := (&models.Strategy{
					RelStrategy:  models.RelStrategy{StrategyID: idg.MustGenSFID()},
					RelProject:   models.RelProject{ProjectID: projectID},
					RelApplet:    models.RelApplet{AppletID: r.Strategies[i].AppletID},
					StrategyInfo: models.StrategyInfo{EventType: r.Strategies[i].EventType, Handler: r.Strategies[i].Handler},
				}).Create(db); err != nil {
					return err
				}
			}
			return nil
		},
	).Do()

	if err != nil {
		l.Error(err)
		return status.CheckDatabaseError(err, "CreateStrategy")
	}

	return
}

func UpdateStrategy(ctx context.Context, strategyID types.SFID, r *CreateStrategyReq) (err error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := models.Strategy{RelStrategy: models.RelStrategy{StrategyID: strategyID}}

	_, l = l.Start(ctx, "UpdateStrategy")
	defer l.End()

	err = sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			return m.FetchByStrategyID(d)
		},
		func(db sqlx.DBExecutor) error {
			m.RelApplet = r.RelApplet
			m.StrategyInfo.EventType = r.EventType
			m.StrategyInfo.Handler = r.Handler
			return m.UpdateByStrategyID(d)
		},
	).Do()

	if err != nil {
		l.Error(err)
		return status.CheckDatabaseError(err, "UpdateStrategy")
	}

	return
}

func GetStrategyByStrategyID(ctx context.Context, strategyID types.SFID) (*models.Strategy, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := models.Strategy{RelStrategy: models.RelStrategy{StrategyID: strategyID}}

	err := m.FetchByStrategyID(d)
	if err != nil {
		return nil, status.CheckDatabaseError(err, "FetchByStrategyID")
	}

	return &m, nil
}

type ListStrategyReq struct {
	projectID   types.SFID
	IDs         []uint64     `in:"query" name:"id,omitempty"`
	AppletIDs   []types.SFID `in:"query" name:"appletID,omitempty"`
	StrategyIDs []types.SFID `in:"query" name:"strategyID,omitempty"`
	EventTypes  []string     `in:"query" name:"eventType,omitempty"`
	datatypes.Pager
}

func (r *ListStrategyReq) SetCurrentProjectID(projectID types.SFID) {
	r.projectID = projectID
}
func (r *ListStrategyReq) Condition() builder.SqlCondition {
	var (
		m  = &models.Strategy{}
		cs []builder.SqlCondition
	)

	cs = append(cs, m.ColProjectID().Eq(r.projectID))
	if len(r.IDs) > 0 {
		cs = append(cs, m.ColID().In(r.IDs))
	}
	if len(r.AppletIDs) > 0 {
		cs = append(cs, m.ColAppletID().In(r.AppletIDs))
	}
	if len(r.StrategyIDs) > 0 {
		cs = append(cs, m.ColStrategyID().In(r.StrategyIDs))
	}
	if len(r.EventTypes) > 0 {
		cs = append(cs, m.ColEventType().In(r.EventTypes))
	}

	return builder.And(cs...)
}

func (r *ListStrategyReq) Additions() builder.Additions {
	m := &models.Strategy{}
	return builder.Additions{
		builder.OrderBy(builder.DescOrder(m.ColCreatedAt())),
		r.Pager.Addition(),
	}
}

type ListStrategyRsp struct {
	Data  []Detail `json:"data"`  // Data strategy data list
	Total int64    `json:"total"` // Total strategy count under current projectID
}

type Detail struct {
	ProjectID  types.SFID   `json:"projectID"`
	Strategies []InfoDetail `json:"strategies,omitempty"`
	datatypes.OperationTimes
}

type InfoDetail struct {
	StrategyID types.SFID `json:"strategyID"`
	AppletID   types.SFID `json:"appletID"`
	AppletName string     `json:"appletName"`
	EventType  string     `json:"eventType"`
	Handler    string     `json:"handler"`
}

type detail struct {
	StrategyID types.SFID `db:"f_strategy_id"`
	AppletID   types.SFID `db:"f_applet_id"`
	AppletName string     `db:"f_applet_name"`
	EventType  string     `db:"f_event_type"`
	Handler    string     `db:"f_handler"`
	datatypes.OperationTimes
}

func ListStrategy(ctx context.Context, r *ListStrategyReq) (*ListStrategyRsp, error) {
	var (
		d    = types.MustMgrDBExecutorFromContext(ctx)
		ret  = &ListStrategyRsp{}
		err  error
		cond = r.Condition()

		mApplet   = &models.Applet{}
		mStrategy = &models.Strategy{}
	)
	ret.Total, err = mStrategy.Count(d, cond)
	if err != nil {
		return nil, status.CheckDatabaseError(err, "CountStrategy")
	}

	details := make([]detail, 0)

	// TODO eventType:applet => 1:n
	err = d.QueryAndScan(
		builder.Select(
			builder.MultiWith(
				",",
				builder.Alias(mStrategy.ColStrategyID(), "f_strategy_id"),
				builder.Alias(mStrategy.ColAppletID(), "f_applet_id"),
				builder.Alias(mApplet.ColName(), "f_applet_name"),
				builder.Alias(mStrategy.ColEventType(), "f_event_type"),
				builder.Alias(mStrategy.ColHandler(), "f_handler"),
				builder.Alias(mStrategy.ColCreatedAt(), "f_created_at"),
				builder.Alias(mStrategy.ColUpdatedAt(), "f_updated_at"),
			),
		).From(
			d.T(mStrategy),
			builder.LeftJoin(d.T(mApplet)).
				On(mStrategy.ColAppletID().Eq(mApplet.ColAppletID())),
			builder.Where(cond),
			builder.OrderBy(
				builder.DescOrder(mStrategy.ColCreatedAt()),
				builder.AscOrder(mApplet.ColName()),
			),
			r.Pager.Addition(),
		),
		&details,
	)
	if err != nil {
		return nil, status.CheckDatabaseError(err, "ListStrategy")
	}

	detailsMap := make(map[types.SFID][]*detail)
	for i := range details {
		appletID := details[i].AppletID
		detailsMap[appletID] = append(detailsMap[appletID], &details[i])
	}

	for _, vmap := range detailsMap {
		infoDetails := make([]InfoDetail, 0, len(vmap))
		for _, v := range vmap {
			if v.AppletID == 0 {
				continue
			}
			infoDetails = append(infoDetails, InfoDetail{
				StrategyID: v.StrategyID,
				AppletID:   v.AppletID,
				AppletName: v.AppletName,
				EventType:  v.EventType,
				Handler:    v.Handler,
			})
		}
		if len(infoDetails) == 0 {
			infoDetails = nil
		}
		ret.Data = append(ret.Data, Detail{
			ProjectID:  r.projectID,
			Strategies: infoDetails,
			OperationTimes: datatypes.OperationTimes{
				CreatedAt: vmap[0].CreatedAt,
				UpdatedAt: vmap[0].UpdatedAt,
			},
		})
	}

	return ret, nil
}

type RemoveStrategyReq struct {
	ProjectName string       `in:"path" name:"projectName"`
	StrategyIDs []types.SFID `in:"query" name:"strategyID"`
}

func RemoveStrategy(ctx context.Context, r *RemoveStrategyReq) error {
	var (
		d         = types.MustMgrDBExecutorFromContext(ctx)
		mStrategy = &models.Strategy{}
		err       error
	)

	return sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			for _, id := range r.StrategyIDs {
				mStrategy.StrategyID = id
				if err = mStrategy.DeleteByStrategyID(d); err != nil {
					return status.CheckDatabaseError(err, "DeleteByStrategyID")
				}
			}
			return nil
		},
	).Do()
}
