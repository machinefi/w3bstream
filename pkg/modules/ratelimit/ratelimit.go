package ratelimit

import (
	"context"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/job"
	"github.com/machinefi/w3bstream/pkg/types"
)

func GetBySFID(ctx context.Context, id types.SFID) (*models.TrafficRateLimit, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.TrafficRateLimit{RelRateLimit: models.RelRateLimit{RateLimitID: id}}

	if err := m.FetchByRateLimitID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.TrafficRateLimitNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func Create(ctx context.Context, r *CreateReq) (*models.TrafficRateLimit, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)
	project := types.MustProjectFromContext(ctx)

	m := &models.TrafficRateLimit{
		RelRateLimit: models.RelRateLimit{RateLimitID: idg.MustGenSFID()},
		RelProject:   models.RelProject{ProjectID: project.ProjectID},
		RateLimitInfo: models.RateLimitInfo{
			Threshold: r.Threshold,
			CycleNum:  r.CycleNum,
			CycleUnit: r.CycleUnit,
			ApiType:   r.ApiType,
		},
	}
	if err := m.Create(d); err != nil {
		if sqlx.DBErr(err).IsConflict() {
			return nil, status.RateLimitConflict
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}

	t := job.NewTrafficTaskWithPrjName(project.Name, *m)
	job.Dispatch(ctx, t)

	return m, nil
}

func Update(ctx context.Context, r *UpdateReq) (*models.TrafficRateLimit, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	project := types.MustProjectFromContext(ctx)
	m := &models.TrafficRateLimit{RelRateLimit: models.RelRateLimit{RateLimitID: r.RateLimitID}}

	err := sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			ctx := types.WithMgrDBExecutor(ctx, d)
			var err error
			m, err = GetBySFID(ctx, r.RateLimitID)
			return err
		},
		func(d sqlx.DBExecutor) error {
			m.RateLimitInfo.Threshold = r.Threshold
			m.RateLimitInfo.CycleNum = r.CycleNum
			m.RateLimitInfo.CycleUnit = r.CycleUnit
			m.RateLimitInfo.ApiType = r.ApiType
			if err := m.UpdateByRateLimitID(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.RateLimitConflict
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
	).Do()

	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}

	t := job.NewTrafficTaskWithPrjName(project.Name, *m)
	job.Dispatch(ctx, t)

	return m, nil
}

func ListByCond(ctx context.Context, r *CondArgs) (data []models.TrafficRateLimit, err error) {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.TrafficRateLimit{}
	)
	data, err = m.List(d, r.Condition())
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return data, nil
}

func ListDetail(ctx context.Context, r *ListReq) (*ListDetailRsp, error) {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)

		rate = &models.TrafficRateLimit{}
		prj  = types.MustProjectFromContext(ctx)
		ret  = &ListDetailRsp{}
		err  error

		cond = r.Condition()
		adds = r.Additions()
	)

	expr := builder.Select(builder.MultiWith(",",
		builder.Alias(prj.ColName(), "f_project_name"),
		rate.ColProjectID(),
		rate.ColRateLimitID(),
		rate.ColThreshold(),
		rate.ColCycleNum(),
		rate.ColCycleUnit(),
		rate.ColApiType(),
		rate.ColCreatedAt(),
		rate.ColUpdatedAt(),
	)).From(
		d.T(rate),
		append([]builder.Addition{
			builder.LeftJoin(d.T(prj)).On(rate.ColProjectID().Eq(prj.ColProjectID())),
			builder.Where(builder.And(cond, prj.ColDeletedAt().Neq(0))),
		}, adds...)...,
	)
	err = d.QueryAndScan(expr, ret.Data)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	ret.Total, err = rate.Count(d, cond)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return ret, nil
}

func GetByProjectAndType(ctx context.Context, id types.SFID, apiType enums.RateLimitApiType) (*models.TrafficRateLimit, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.TrafficRateLimit{
		RelProject:    models.RelProject{ProjectID: id},
		RateLimitInfo: models.RateLimitInfo{ApiType: apiType},
	}

	if err := m.FetchByProjectIDAndApiType(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.TrafficRateLimitNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}
