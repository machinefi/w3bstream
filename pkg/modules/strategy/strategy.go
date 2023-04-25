package strategy

import (
	"context"
	"time"

	"github.com/pkg/errors"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

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

func Update(ctx context.Context, id types.SFID, r *UpdateReq) (err error) {
	var m *models.Strategy

	err = sqlx.NewTasks(types.MustMgrDBExecutorFromContext(ctx)).With(
		func(d sqlx.DBExecutor) error {
			m, err = GetBySFID(ctx, id)
			return err
		},
		func(d sqlx.DBExecutor) error {
			m.RelApplet, m.StrategyInfo = r.RelApplet, r.StrategyInfo
			if err = m.UpdateByStrategyID(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.StrategyConflict
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
	).Do()
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

func Remove(ctx context.Context, r *CondArgs) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Strategy{}

	prj := types.MustProjectFromContext(ctx)

	_, err := d.Exec(
		builder.Update(d.T(m)).Set(
			m.ColDeletedAt().ValueBy(time.Now().Unix()),
		).Where(r.Condition(prj.ProjectID)),
	)
	if err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}

	return nil
}

func List(ctx context.Context, r *ListReq) (ret *ListRsp, err error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Strategy{}

	prj := types.MustProjectFromContext(ctx)
	cond := r.Condition(prj.ProjectID)
	adds := builder.Additions{
		r.Pager.Addition(),
		builder.OrderBy(builder.DescOrder(m.ColUpdatedAt())),
		builder.OrderBy(builder.DescOrder(m.ColCreatedAt())),
	}

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

func Create(ctx context.Context, r *CreateReq) (*CreateRsp, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)

	if len(r.Data) == 0 {
		return &CreateRsp{}, nil
	}

	prj := types.MustProjectFromContext(ctx).ProjectID
	ids := confid.MustSFIDGeneratorFromContext(ctx).MustGenSFIDs(len(r.Data))
	ret := &CreateRsp{Data: make([]*models.Strategy, 0, len(r.Data))}

	err := sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			for i := range r.Data {
				m := &models.Strategy{
					RelStrategy:  models.RelStrategy{StrategyID: ids[i]},
					RelProject:   models.RelProject{ProjectID: prj},
					RelApplet:    r.Data[i].RelApplet,
					StrategyInfo: r.Data[i].StrategyInfo,
				}
				if err := m.Create(d); err != nil {
					if sqlx.DBErr(err).IsConflict() {
						// TODO gen model.MayBeConflictFields() for more hint to frontend
						return status.StrategyConflict
					}
					return status.DatabaseError.StatusErr().WithDesc(err.Error())
				}
				ret.Data = append(ret.Data, m)
			}
			return nil
		},
	).Do()

	if err != nil {
		return nil, err
	}
	return ret, nil
}
