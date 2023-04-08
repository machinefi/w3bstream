package ratelimit

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm/kvdb"
)

type CreateTrafficRateLimitReq struct {
	models.RateLimitInfo
}

func CreateRateLimit(ctx context.Context, r *CreateTrafficRateLimitReq) (*models.TrafficRateLimit, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)
	project := types.MustProjectFromContext(ctx)

	_, l = l.Start(ctx, "CreateTrafficRateLimit")
	defer l.End()

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
		l.Error(err)
		return nil, err
	}

	//t := job.NewTrafficTask(*m)
	//job.Dispatch(ctx, t)
	//time.Sleep(3 * time.Second)

	// TODO delete this, use TrafficTask.Scheduler
	rDB := kvdb.MustRedisDBKeyFromContext(ctx)
	prj := types.MustProjectFromContext(ctx)
	schedulerJobs := types.MustSchedulerJobsFromContext(ctx)
	s, ok := schedulerJobs.Jobs.Load(prj.Name)
	if !ok || s == nil {
		s := gocron.NewScheduler(time.UTC)
		genSchedulerJob(prj.Name, m.RateLimitInfo, rDB, s)
		s.StartImmediately()
		s.StartAsync()
		schedulerJobs.Jobs.Store(prj.Name, s)
	} else {
		s.Clear()
		genSchedulerJob(prj.Name, m.RateLimitInfo, rDB, s)
	}

	return m, nil
}

func UpdateRateLimit(ctx context.Context, rateLimitID types.SFID, r *CreateTrafficRateLimitReq) (*models.TrafficRateLimit, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.TrafficRateLimit{RelRateLimit: models.RelRateLimit{RateLimitID: rateLimitID}}

	_, l = l.Start(ctx, "UpdateTrafficRateLimit")
	defer l.End()

	err := sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			return m.FetchByRateLimitID(d)
		},
		func(db sqlx.DBExecutor) error {
			m.RateLimitInfo.Threshold = r.Threshold
			m.RateLimitInfo.CycleNum = r.CycleNum
			m.RateLimitInfo.CycleUnit = r.CycleUnit
			m.RateLimitInfo.ApiType = r.ApiType
			return m.UpdateByRateLimitID(d)
		},
	).Do()

	if err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "UpdateTrafficRateLimit")
	}

	// TODO update window

	return m, nil
}

func genSchedulerJob(projectName string, rateLimitInfo models.RateLimitInfo, rDB *kvdb.RedisDB, s *gocron.Scheduler) {
	now := time.Now().UTC()
	seconds := 0
	switch rateLimitInfo.CycleUnit {
	case enums.TRAFFIC_CYCLE__MINUTE:
		nextMinute := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
		s.Every(rateLimitInfo.CycleNum).Minutes().StartAt(nextMinute)
		seconds = rateLimitInfo.CycleNum * 60
	case enums.TRAFFIC_CYCLE__HOUR:
		nextHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
		s.Every(rateLimitInfo.CycleNum).Hours().StartAt(nextHour)
		seconds = rateLimitInfo.CycleNum * 60 * 60
	case enums.TRAFFIC_CYCLE__DAY:
		nextDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		s.Every(rateLimitInfo.CycleNum).Day().StartAt(nextDay)
		seconds = rateLimitInfo.CycleNum * 60 * 60 * 24
	case enums.TRAFFIC_CYCLE__MONTH:
		nextMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		s.Every(rateLimitInfo.CycleNum).Day().StartAt(nextMonth)
		seconds = rateLimitInfo.CycleNum * 60 * 60 * 24 * 31
	case enums.TRAFFIC_CYCLE__YEAR:
		nextYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		s.Every(rateLimitInfo.CycleNum).Day().StartAt(nextYear)
		seconds = rateLimitInfo.CycleNum * 60 * 60 * 24 * 31 * 366
	}
	s.Do(resetWindow, projectName, rateLimitInfo.Threshold, int64(seconds), rDB)
}

func resetWindow(projectName string, threshold int, exp int64, rDB *kvdb.RedisDB) {
	err := rDB.SetKeyWithEX(projectName, []byte(strconv.Itoa(threshold)), exp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(strconv.Itoa(threshold) + "s" + time.Now().Format("2006-01-02 15:04:05"))
}
