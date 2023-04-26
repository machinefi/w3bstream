package cronjob

import (
	"context"

	"github.com/robfig/cron/v3"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateReq struct {
	ProjectID types.SFID `json:"-"`
	models.CronJobInfo
}

type CondArgs struct {
	ProjectID  types.SFID   `name:"-"`
	CronJobIDs []types.SFID `in:"query" name:"cronJobID,omitempty"`
	EventTypes []string     `in:"query" name:"eventType,omitempty"`
}

func (r *CondArgs) Condition() builder.SqlCondition {
	var (
		m  = &models.CronJob{}
		cs []builder.SqlCondition
	)

	if r.ProjectID != 0 {
		cs = append(cs, m.ColProjectID().Eq(r.ProjectID))
	}
	if len(r.CronJobIDs) > 0 {
		cs = append(cs, m.ColCronJobID().In(r.CronJobIDs))
	}
	if len(r.EventTypes) > 0 {
		cs = append(cs, m.ColEventType().In(r.EventTypes))
	}
	cs = append(cs, m.ColDeletedAt().Eq(0))
	return builder.And(cs...)
}

type ListReq struct {
	CondArgs
	datatypes.Pager
}

func (r *ListReq) Additions() builder.Additions {
	m := &models.CronJob{}
	return builder.Additions{
		builder.OrderBy(
			builder.DescOrder(m.ColUpdatedAt()),
			builder.DescOrder(m.ColCreatedAt()),
		),
		r.Pager.Addition(),
	}
}

type ListRsp struct {
	Data  []models.CronJob `json:"data"`
	Total int64            `json:"total"`
}

func List(ctx context.Context, r *ListReq) (*ListRsp, error) {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.CronJob{}

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

func Create(ctx context.Context, r *CreateReq) (*models.CronJob, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	if _, err := cron.ParseStandard(r.CronExpressions); err != nil {
		return nil, status.InvalidCronExpressions.StatusErr().WithDesc(err.Error())
	}

	n := *r
	n.EventType = getEventType(n.EventType)
	m := &models.CronJob{
		RelCronJob:  models.RelCronJob{CronJobID: idg.MustGenSFID()},
		RelProject:  models.RelProject{ProjectID: r.ProjectID},
		CronJobInfo: n.CronJobInfo,
	}
	if err := m.Create(d); err != nil {
		if sqlx.DBErr(err).IsConflict() {
			return nil, status.CronJobConflict
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func RemoveBySFID(ctx context.Context, id types.SFID) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.CronJob{RelCronJob: models.RelCronJob{CronJobID: id}}

	if err := m.DeleteByCronJobID(d); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func GetBySFID(ctx context.Context, id types.SFID) (*models.CronJob, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)

	m := &models.CronJob{RelCronJob: models.RelCronJob{CronJobID: id}}
	if err := m.FetchByCronJobID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.CronJobNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func getEventType(eventType string) string {
	if eventType == "" {
		return enums.EVENTTYPEDEFAULT
	}
	return eventType
}
