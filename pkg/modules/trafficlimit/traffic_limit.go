package trafficlimit

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

func GetBySFID(ctx context.Context, id types.SFID) (*models.TrafficLimit, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.TrafficLimit{RelTrafficLimit: models.RelTrafficLimit{TrafficLimitID: id}}

	if err := m.FetchByTrafficLimitID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.TrafficLimitNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func Create(ctx context.Context, r *CreateReq) (*models.TrafficLimit, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)
	project := types.MustProjectFromContext(ctx)

	m := &models.TrafficLimit{
		RelTrafficLimit: models.RelTrafficLimit{TrafficLimitID: idg.MustGenSFID()},
		RelProject:      models.RelProject{ProjectID: project.ProjectID},
		TrafficLimitInfo: models.TrafficLimitInfo{
			Threshold: r.Threshold,
			Duration:  r.Duration,
			ApiType:   r.ApiType,
		},
	}

	err := sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			if err := m.Create(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.TrafficLimitConflict
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			if err := CreateScheduler(ctx, fmt.Sprintf("%s::%s", project.Name, r.ApiType.String()), &m.TrafficLimitInfo); err != nil {
				return status.CreateTrafficSchedulerFailed
			}
			return nil
		},
	).Do()
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}

	return m, nil
}

func Update(ctx context.Context, r *UpdateReq) (*models.TrafficLimit, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	project := types.MustProjectFromContext(ctx)
	m := &models.TrafficLimit{RelTrafficLimit: models.RelTrafficLimit{TrafficLimitID: r.TrafficLimitID}}

	err := sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			ctx := types.WithMgrDBExecutor(ctx, d)
			var err error
			m, err = GetBySFID(ctx, r.TrafficLimitID)
			return err
		},
		func(d sqlx.DBExecutor) error {
			m.TrafficLimitInfo.Threshold = r.Threshold
			m.TrafficLimitInfo.Duration = r.Duration
			m.TrafficLimitInfo.ApiType = r.ApiType
			if err := m.UpdateByTrafficLimitID(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.TrafficLimitConflict
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			if err := UpdateScheduler(ctx, fmt.Sprintf("%s::%s", project.Name, r.ApiType.String()), &m.TrafficLimitInfo); err != nil {
				return status.UpdateTrafficSchedulerFailed
			}
			return nil
		},
	).Do()
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}

	return m, nil
}

func List(ctx context.Context, r *ListReq) (*ListRsp, error) {
	var (
		d       = types.MustMgrDBExecutorFromContext(ctx)
		traffic = &models.TrafficLimit{}
		ret     = &ListRsp{}
		cond    = r.Condition()

		err error
	)

	if ret.Data, err = traffic.List(d, cond, r.Addition()); err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	if ret.Total, err = traffic.Count(d, cond); err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return ret, nil
}

func ListByCond(ctx context.Context, r *CondArgs) (data []models.TrafficLimit, err error) {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.TrafficLimit{}
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

		rate = &models.TrafficLimit{}
		prj  = types.MustProjectFromContext(ctx)
		ret  = &ListDetailRsp{}
		err  error

		cond = r.Condition()
		adds = r.Additions()
	)

	expr := builder.Select(builder.MultiWith(",",
		builder.Alias(prj.ColName(), "f_project_name"),
		rate.ColProjectID(),
		rate.ColTrafficLimitID(),
		rate.ColThreshold(),
		rate.ColDuration(),
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

func GetByProjectAndType(ctx context.Context, id types.SFID, apiType enums.TrafficLimitType) (*models.TrafficLimit, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.TrafficLimit{
		RelProject:       models.RelProject{ProjectID: id},
		TrafficLimitInfo: models.TrafficLimitInfo{ApiType: apiType},
	}

	if err := m.FetchByProjectIDAndApiType(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.TrafficLimitNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}
