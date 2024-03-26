package trafficlimit

import (
	"context"
	"time"

	"github.com/go-co-op/gocron/v2"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

var prefix = "traffic_limit"

func Init(ctx context.Context) error {
	_, l := logr.Start(ctx, "trafficLimit.Init")
	defer l.End()

	kv := types.MustRedisEndpointFromContext(ctx).WithPrefix(prefix)
	d := types.MustMgrDBExecutorFromContext(ctx)

	if !types.EnableTrafficLimitFromContext(ctx) {
		return nil
	}

	keys, _ := kv.Keys("*")
	_ = kv.RawDel(keys...)

	s, err := gocron.NewScheduler()
	if err != nil {
		return err
	}
	_, err = s.NewJob(gocron.DurationJob(time.Minute), gocron.NewTask(func(d sqlx.DBExecutor) error {
		list, err := (&models.TrafficLimit{}).List(d, nil)
		if err != nil {
			return err
		}

		ids := make(map[types.SFID]int)
		for i := range list {
			ids[list[i].TrafficLimitID] = i
		}

		for _, i := range ids {
			_ = AddAndStartScheduler(ctx, &list[i])
		}

		keys := make(map[types.SFID]struct{})
		schedulers.Range(func(k types.SFID, v *Scheduler) bool {
			keys[k] = struct{}{}
			return true
		})

		for id, _ := range keys {
			if _, ok := ids[id]; !ok {
				RmvScheduler(ctx, id)
			}
		}

		return nil
	}, d))
	if err != nil {
		return err
	}
	s.Start()
	return nil
}

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
	var (
		d   = types.MustMgrDBExecutorFromContext(ctx)
		idg = confid.MustSFIDGeneratorFromContext(ctx)
		prj = types.MustProjectFromContext(ctx)
	)

	m := &models.TrafficLimit{
		RelTrafficLimit: models.RelTrafficLimit{TrafficLimitID: idg.MustGenSFID()},
		RelProject:      models.RelProject{ProjectID: prj.ProjectID},
		TrafficLimitInfo: models.TrafficLimitInfo{
			Threshold: r.Threshold,
			Duration:  r.Duration,
			ApiType:   r.ApiType,
			StartAt:   types.Timestamp{Time: time.Now()},
		},
	}

	if err := m.Create(d); err != nil {
		if sqlx.DBErr(err).IsConflict() {
			return nil, status.TrafficLimitConflict
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func Update(ctx context.Context, id types.SFID, r *UpdateReq) (*models.TrafficLimit, error) {
	var (
		d   = types.MustMgrDBExecutorFromContext(ctx)
		m   *models.TrafficLimit
		err error
	)

	err = sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			m, err = GetBySFID(ctx, id)
			if err != nil {
				return err
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			m.TrafficLimitID = id
			m.Threshold = r.Threshold
			m.Duration = r.Duration
			if err = m.UpdateByTrafficLimitID(d); err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
	).Do()
	if err != nil {
		return nil, err
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
		builder.Alias(rate.ColProjectID(), "f_project_id"),
		builder.Alias(rate.ColTrafficLimitID(), "f_traffic_limit_id"),
		builder.Alias(rate.ColThreshold(), "f_threshold"),
		builder.Alias(rate.ColDuration(), "f_duration"),
		builder.Alias(rate.ColApiType(), "f_api_type"),
		builder.Alias(rate.ColCreatedAt(), "f_created_at"),
		builder.Alias(rate.ColUpdatedAt(), "f_updated_at"),
	)).From(
		d.T(rate),
		append([]builder.Addition{
			builder.LeftJoin(d.T(prj)).On(rate.ColProjectID().Eq(prj.ColProjectID())),
			builder.Where(cond),
		}, adds...)...,
	)
	err = d.QueryAndScan(expr, &ret.Data)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	ret.Total, err = rate.Count(d, cond)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return ret, nil
}

func RemoveBySFID(ctx context.Context, id types.SFID) error {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.TrafficLimit{}
	)

	m.TrafficLimitID = id
	if err := m.DeleteByTrafficLimitID(d); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func Remove(ctx context.Context, r *CondArgs) error {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.TrafficLimit{}

		err error
	)

	if r.Condition().IsNil() {
		return status.InvalidDeleteCondition
	}

	return sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			_, err = ListDetail(ctx, &ListReq{CondArgs: *r})
			if err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			expr := builder.Delete().From(d.T(m), builder.Where(r.Condition()))
			_, err = d.Exec(expr)
			if err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
	).Do()
}

func TrafficLimit(ctx context.Context, prj types.SFID, tpe enums.TrafficLimitType) error {
	ctx, l := logr.Start(ctx, "trafficLimit.TrafficLimit")
	defer l.End()

	r := types.MustRedisEndpointFromContext(ctx)

	// for statistics per hour
	stat := r.WithPrefix("stat" + ":" + prj.String())
	if total, err := stat.IncrBy(time.Now().Format("2006010215"), 1); err != nil {
		l.WithValues("prj", prj, "total", total).Info("")
	}

	// for traffic limit
	limit := r.WithPrefix(prefix)
	l = l.WithValues("prj", prj, "tpe", tpe)
	m := &models.TrafficLimit{
		RelProject:       models.RelProject{ProjectID: prj},
		TrafficLimitInfo: models.TrafficLimitInfo{ApiType: tpe},
	}

	exists, _ := limit.Exists(m.CacheKey())
	if !exists {
		l.Info("no strategy")
		return nil
	}

	count, _ := limit.IncrBy(m.CacheKey(), -1)
	if count <= 0 {
		l.Info("limited")
		return status.TrafficLimitExceededFailed
	}

	l.WithValues("remain", count).Info("")
	return nil
}
