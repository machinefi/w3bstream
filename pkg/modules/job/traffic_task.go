package job

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm/kvdb"
)

func NewTrafficTaskWithPrjKey(projectKey string, traffic models.TrafficLimit) *TrafficTask {
	return &TrafficTask{projectKey: projectKey, traffic: traffic}
}

type TrafficTask struct {
	projectKey string
	traffic    models.TrafficLimit
	mq.TaskState
}

var _ mq.Task = (*TrafficTask)(nil)

func (t *TrafficTask) Subject() string {
	return "TrafficLimitTask"
}

func (t *TrafficTask) Arg() interface{} {
	return t
}

func (t *TrafficTask) ID() string {
	return fmt.Sprintf("%s::%s", t.Subject(), t.traffic.TrafficLimitID)
}

func (t *TrafficTask) Scheduler(ctx context.Context) {
	rDB := kvdb.MustRedisDBKeyFromContext(ctx)
	schedulerJobs := types.MustSchedulerJobsFromContext(ctx)

	s, ok := schedulerJobs.Jobs.Load(t.projectKey)
	if ok && s != nil {
		s.Clear()
	}
	s = gocron.NewScheduler(time.UTC)
	genSchedulerJob(t.projectKey, t.traffic.TrafficLimitInfo, rDB, s)
	s.StartImmediately()
	s.StartAsync()
	schedulerJobs.Jobs.Store(t.projectKey, s)
}

func genSchedulerJob(projectKey string, rateLimitInfo models.TrafficLimitInfo, rDB *kvdb.RedisDB, s *gocron.Scheduler) {
	now := time.Now().UTC()
	seconds := rateLimitInfo.Duration.Duration().Seconds()
	if seconds >= 24*time.Hour.Seconds() {
		nextDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		s.Every(seconds).Second().StartAt(nextDay)
	} else if seconds >= time.Hour.Seconds() {
		nextHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
		s.Every(seconds).Second().StartAt(nextHour)
	} else if seconds >= time.Minute.Seconds() {
		nextMinute := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
		s.Every(seconds).Second().StartAt(nextMinute)
	} else {
		nextSecond := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, now.Location())
		s.Every(seconds).Second().StartAt(nextSecond)
	}

	//switch rateLimitInfo.CycleUnit {
	//case enums.TRAFFIC_CYCLE__MINUTE:
	//	nextMinute := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	//	s.Every(rateLimitInfo.CycleNum).Minutes().StartAt(nextMinute)
	//	seconds = rateLimitInfo.CycleNum * 60
	//case enums.TRAFFIC_CYCLE__HOUR:
	//	nextHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	//	s.Every(rateLimitInfo.CycleNum).Hours().StartAt(nextHour)
	//	seconds = rateLimitInfo.CycleNum * 60 * 60
	//case enums.TRAFFIC_CYCLE__DAY:
	//	nextDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	//	s.Every(rateLimitInfo.CycleNum).Day().StartAt(nextDay)
	//	seconds = rateLimitInfo.CycleNum * 60 * 60 * 24
	//case enums.TRAFFIC_CYCLE__MONTH:
	//	nextMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	//	s.Every(rateLimitInfo.CycleNum).Day().StartAt(nextMonth)
	//	seconds = rateLimitInfo.CycleNum * 60 * 60 * 24 * 31
	//case enums.TRAFFIC_CYCLE__YEAR:
	//	nextYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	//	s.Every(rateLimitInfo.CycleNum).Day().StartAt(nextYear)
	//	seconds = rateLimitInfo.CycleNum * 60 * 60 * 24 * 31 * 366
	//}
	s.Do(resetWindow, projectKey, rateLimitInfo.Threshold, int64(seconds), rDB)
}

// TODO redeploy reset redis yue
func resetWindow(projectKey string, threshold int, exp int64, rDB *kvdb.RedisDB) error {
	err := rDB.SetKeyWithEX(projectKey, []byte(strconv.Itoa(threshold)), exp)
	if err != nil {
		return err
	}
	// TODO del
	fmt.Println(projectKey + " - " + strconv.Itoa(threshold) + "s" + time.Now().Format("2006-01-02 15:04:05"))
	return nil
}
