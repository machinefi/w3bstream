package tasks

import (
	"context"
	"reflect"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/modules/job"
)

type RateLimitTask struct {
	SchedulerJob *job.TrafficTask
}

func (s *RateLimitTask) SetArg(v interface{}) error {
	if job, ok := v.(*job.TrafficTask); ok {
		s.SchedulerJob = job
		return nil
	}
	return errors.Errorf("invalid arg: %s", reflect.TypeOf(v))
}

func (s *RateLimitTask) Output(ctx context.Context) (interface{}, error) {
	s.SchedulerJob.Scheduler(ctx)
	return nil, nil
}
