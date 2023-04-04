package job

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
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
	// TODO project is null
	prj := types.MustProjectFromContext(ctx)
	schedulerJobs := types.MustSchedulerJobsFromContext(ctx)
	s, ok := schedulerJobs.Jobs.Load(prj.Name)
	if !ok || s == nil {
		s := gocron.NewScheduler(time.UTC)
		//s.Every(t.traffic.RateLimitInfo.Duration).Seconds().Do(resetWindow, t.traffic.RateLimitInfo.Count)
		s.Every(10).Seconds().Do(resetWindow, 10)
		s.StartAsync()
		schedulerJobs.Jobs.Store(prj.Name, s)
	} else {
		s.Clear()
		s.Every(20).Seconds().Do(resetWindow, 20)
	}
}

func resetWindow(count int) {
	// set redis count
	fmt.Println(string(count) + "s" + time.Now().Format("2006-01-02 15:04:05"))
}
