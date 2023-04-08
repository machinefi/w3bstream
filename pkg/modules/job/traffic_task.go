package job

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm/kvdb"
)

func NewTrafficTask(traffic models.TrafficRateLimit) *TrafficTask {
	return &TrafficTask{traffic: traffic}
}

type TrafficTask struct {
	traffic models.TrafficRateLimit
	mq.TaskState
}

var _ mq.Task = (*TrafficTask)(nil)

func (t *TrafficTask) Subject() string {
	return "RateLimitTask"
}

func (t *TrafficTask) Arg() interface{} {
	return t
}

func (t *TrafficTask) ID() string {
	return fmt.Sprintf("%s::%s", t.Subject(), t.traffic.RateLimitID)
}

func (t *TrafficTask) Scheduler(ctx context.Context) {
	rDB := kvdb.MustRedisDBKeyFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	// TODO project is null
	prj := types.MustProjectFromContext(ctx)
	schedulerJobs := types.MustSchedulerJobsFromContext(ctx)

	_, l = l.Start(ctx, "trafficTaskScheduler")
	defer l.End()

	s, ok := schedulerJobs.Jobs.Load(prj.Name)
	if !ok || s == nil {
		s := gocron.NewScheduler(time.UTC)
		genSchedulerJob(prj.Name, t.traffic.RateLimitInfo, l, rDB, s)
		s.StartImmediately()
		s.StartAsync()
		schedulerJobs.Jobs.Store(prj.Name, s)
	} else {
		s.Clear()
		genSchedulerJob(prj.Name, t.traffic.RateLimitInfo, l, rDB, s)
	}
}

func genSchedulerJob(projectName string, rateLimitInfo models.RateLimitInfo, l log.Logger, rDB *kvdb.RedisDB, s *gocron.Scheduler) {
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
	s.Do(resetWindow, projectName, rateLimitInfo.Threshold, seconds, l, rDB)
}

func resetWindow(projectName string, threshold int, exp int64, l log.Logger, rDB *kvdb.RedisDB) {
	err := rDB.SetKeyWithEX(projectName, []byte(strconv.Itoa(threshold)), exp)
	if err != nil {
		l.Error(err)
	}
	fmt.Println(strconv.Itoa(threshold) + "s" + time.Now().Format("2006-01-02 15:04:05"))
}
