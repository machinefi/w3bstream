package trafficlimit

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types/wasm/kvdb"
)

func genScheduler(trafficScheduler *TrafficScheduler) (*TrafficScheduler, error) {
	trafficScheduler.rDB = kvdb.MustRedisDBKeyFromContext(trafficScheduler.ctx)
	trafficScheduler.sch = gocron.NewScheduler(time.UTC)
	trafficScheduler.genSchedulerJob()
	trafficSchedulers.Store(trafficScheduler.projectKey, trafficScheduler)
	if err := trafficScheduler.Do(); err != nil {
		return nil, err
	}
	return trafficScheduler, nil
}

func CreateScheduler(ctx context.Context, projectKey string, trafficInfo *models.TrafficLimitInfo) error {
	ts, err := genScheduler(&TrafficScheduler{
		ctx:         ctx,
		projectKey:  projectKey,
		trafficInfo: trafficInfo,
	})
	if err != nil {
		return err
	}

	ts.StartNow()
	return nil
}

func UpdateScheduler(ctx context.Context, projectKey string, trafficInfo *models.TrafficLimitInfo) error {
	ts, ok := trafficSchedulers.Load(projectKey)
	if ok && ts != nil {
		ts.Stop()
		trafficSchedulers.Remove(projectKey)
	}

	ts, err := genScheduler(&TrafficScheduler{
		ctx:         ctx,
		projectKey:  projectKey,
		trafficInfo: trafficInfo,
	})
	if err != nil {
		return err
	}
	ts.Start()
	return nil
}

func RestartScheduler(ctx context.Context, projectKey string, trafficInfo *models.TrafficLimitInfo) error {
	ts, err := genScheduler(&TrafficScheduler{
		ctx:         ctx,
		projectKey:  projectKey,
		trafficInfo: trafficInfo,
	})
	if err != nil {
		return err
	}
	ts.Start()
	return nil
}

func DeleteScheduler(projectKey string) error {
	ts, ok := trafficSchedulers.Load(projectKey)
	if !ok {
		return errors.New("trafficScheduler not found")
	}
	if ts != nil {
		ts.Stop()
	}
	trafficSchedulers.Remove(projectKey)
	// TODO del redis key
	return nil
}

var trafficSchedulers = *mapx.New[string, *TrafficScheduler]()

type TrafficScheduler struct {
	ctx         context.Context
	projectKey  string
	trafficInfo *models.TrafficLimitInfo
	sch         *gocron.Scheduler
	rDB         *kvdb.RedisDB
}

func NewTrafficScheduler(ctx context.Context, projectKey string, trafficInfo *models.TrafficLimitInfo) *TrafficScheduler {
	ts := &TrafficScheduler{
		projectKey:  projectKey,
		trafficInfo: trafficInfo,
	}
	ts.rDB = kvdb.MustRedisDBKeyFromContext(ctx)
	ts.sch = gocron.NewScheduler(time.UTC)
	ts.genSchedulerJob()
	return ts
}

func (ts *TrafficScheduler) genSchedulerJob() {
	now := time.Now().UTC()
	seconds := ts.trafficInfo.Duration.Duration().Seconds()
	if seconds >= 24*time.Hour.Seconds() {
		nextDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		ts.sch.Every(int(seconds)).Second().StartAt(nextDay)
	} else if seconds >= time.Hour.Seconds() {
		nextHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
		ts.sch.Every(int(seconds)).Second().StartAt(nextHour)
	} else if seconds >= time.Minute.Seconds() {
		nextMinute := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
		ts.sch.Every(int(seconds)).Second().StartAt(nextMinute)
	} else {
		nextSecond := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, now.Location())
		ts.sch.Every(int(seconds)).Second().StartAt(nextSecond)
	}
}

func (ts *TrafficScheduler) StartNow() {
	ts.sch.StartImmediately().StartAsync()
}

func (ts *TrafficScheduler) Start() {
	ts.sch.StartAsync()
}

func (ts *TrafficScheduler) Stop() {
	ts.sch.Stop()
}

func (ts *TrafficScheduler) Do() error {
	_, err := ts.sch.Do(resetWindow, ts.projectKey, ts.trafficInfo.Threshold,
		int64(ts.trafficInfo.Duration.Duration().Seconds()), ts.rDB)
	if err != nil {
		return err
	}
	return nil
}

func resetWindow(projectKey string, threshold int, exp int64, rDB *kvdb.RedisDB) error {
	err := rDB.SetKeyWithEX(projectKey, []byte(strconv.Itoa(threshold)), exp)
	if err != nil {
		return err
	}
	// TODO del
	fmt.Println(projectKey + " - " + strconv.Itoa(threshold) + "s" + time.Now().Format("2006-01-02 15:04:05"))
	return nil
}
