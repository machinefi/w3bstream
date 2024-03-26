package trafficlimit

import (
	"context"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	"github.com/machinefi/w3bstream/pkg/depends/conf/redis"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func AddScheduler(ctx context.Context, v *models.TrafficLimit, start ...bool) error {
	_, l := logger.NewSpanContext(ctx, "traffic_limit.AddScheduler")
	defer l.End()

	l = l.WithValues("traffic_limit", v.TrafficLimitID, "du", v.Duration, "count", v.Threshold)
	s, ok := schedulers.Load(v.TrafficLimitID)
	if !ok {
		s, err := NewScheduler(ctx, v, start...)
		if err != nil {
			l.Error(err)
			return err
		}
		schedulers.Store(v.TrafficLimitID, s)
		l.Info("schedule created")
		return nil
	}
	s.update(ctx, v)
	return nil
}

func AddAndStartScheduler(ctx context.Context, v *models.TrafficLimit) error {
	return AddScheduler(ctx, v, true)
}

func RmvScheduler(ctx context.Context, id types.SFID) {
	_, l := logger.NewSpanContext(ctx, "traffic_limit.RmvScheduler")
	defer l.End()

	s, _ := schedulers.Load(id)
	if s != nil {
		s.Stop()
	}
	schedulers.Remove(id)
	l.WithValues("traffic_limit", id).Info("schedule removed")
}

var schedulers = *mapx.New[types.SFID, *Scheduler]()

func NewScheduler(ctx context.Context, v *models.TrafficLimit, start ...bool) (*Scheduler, error) {
	s := &Scheduler{
		key: v.CacheKey(),
		du:  v.Duration.Duration(),
		cnt: int64(v.Threshold),
		kv:  types.MustRedisEndpointFromContext(ctx).WithPrefix(prefix),
	}
	var err error
	s.sch, err = gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	_, err = s.sch.NewJob(gocron.DurationJob(s.du), gocron.NewTask(s.reset))
	if err != nil {
		return nil, err
	}
	if len(start) > 0 && start[0] {
		s.Start()
	}
	return s, nil
}

type Scheduler struct {
	mu  sync.Mutex
	key string
	du  time.Duration
	cnt int64
	sch gocron.Scheduler
	kv  *redis.Redis
}

func (s *Scheduler) Start() {
	s.reset()
	s.sch.Start()
}

func (s *Scheduler) Stop() {
	_ = s.kv.Del(s.key)
	_ = s.sch.Shutdown()
}

func (s *Scheduler) reset() {
	s.mu.Lock()
	key, cnt, du := s.key, s.cnt, s.du
	s.mu.Unlock()
	_ = s.kv.SetEx(key, cnt, du)
}

func (s *Scheduler) update(ctx context.Context, v *models.TrafficLimit) {
	_, l := logger.NewSpanContext(ctx, "traffic_limit.AddScheduler")
	defer l.End()

	l = l.WithValues("traffic_limit", v.TrafficLimitID)

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.key != v.CacheKey() || s.cnt != int64(v.Threshold) || s.du != v.Duration.Duration() {
		s.key = v.CacheKey()
		s.du = v.Duration.Duration()
		s.cnt = int64(v.Threshold)
		err := s.kv.SetEx(s.key, s.cnt, s.du)
		l = l.WithValues("traffic_limit", v.TrafficLimitInfo, "du", v.Duration, "count", s.cnt)
		if err != nil {
			l.Error(errors.Wrap(err, "schedule update failed: %v"))
		} else {
			l.Info("schedule updated")
		}
	}
}
