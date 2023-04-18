package cronjob

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/types"
)

const (
	listIntervalSecond = 3
)

type cronJob struct {
	listIntervalSecond int
}

func (t *cronJob) run(ctx context.Context) {
	l := types.MustLoggerFromContext(ctx)
	s := gocron.NewScheduler(time.UTC)
	s.TagsUnique()

	if _, err := s.Every(t.listIntervalSecond).Seconds().Do(t.do, ctx, s); err != nil {
		l.Fatal(errors.Wrap(err, "create cronjob main loop failed"))
	}
	s.StartAsync()
}

func (t *cronJob) do(ctx context.Context, s *gocron.Scheduler) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.CronJob{}

	_, l = l.Start(ctx, "cronjob.run")
	defer l.End()

	cs, err := m.List(d, nil)
	if err != nil {
		l.Error(errors.Wrap(err, "list cronjob db failed"))
		return
	}

	t.tidyCronJobs(ctx, s, cs, l)

	for _, c := range cs {
		if _, err := s.Cron(c.CronExpressions).Tag(c.CronJobID.String()).Do(t.sendEvent, ctx, c); err != nil {
			if !strings.Contains(err.Error(), "non-unique tag") {
				l.WithValues("cronJobID", c.CronJobID).Error(errors.Wrap(err, "create new cron job failed"))
			}
		}
	}
}

func (t *cronJob) tidyCronJobs(ctx context.Context, s *gocron.Scheduler, cs []models.CronJob, l log.Logger) {
	cronJobIDs := make(map[types.SFID]bool, len(cs))
	for _, c := range cs {
		cronJobIDs[c.CronJobID] = true
	}
	for _, tag := range s.GetAllTags() {
		id, err := strconv.ParseUint(tag, 10, 64)
		if err != nil {
			l.WithValues("tag", tag).Error(errors.Wrap(err, "parse tag to uint64 failed"))
			continue
		}
		if !cronJobIDs[types.SFID(id)] {
			if err := s.RemoveByTag(tag); err != nil {
				l.WithValues("tag", tag).Error(errors.Wrap(err, "remove cron job failed"))
			} else {
				l.WithValues("tag", tag).Info("remove cron job success")
			}
		}
	}
}

func (t *cronJob) sendEvent(ctx context.Context, c models.CronJob) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "cronjob.sendEvent")
	defer l.End()
	l = l.WithValues("cronJobID", c.CronJobID)

	m := &models.Project{RelProject: models.RelProject{ProjectID: c.ProjectID}}
	if err := m.FetchByProjectID(d); err != nil {
		l.Error(errors.Wrap(err, "get project failed"))
		return
	}

	e := eventpb.Event{
		Header: &eventpb.Header{
			EventType: c.EventType,
		},
		Payload: []byte(fmt.Sprintf("cronJobID:%d", c.CronJobID)),
	}
	if _, err := event.OnEventReceived(ctx, m.ProjectName.Name, &e); err != nil {
		l.Error(errors.Wrap(err, "send event failed"))
	}
}

func Run(ctx context.Context) {
	c := &cronJob{
		listIntervalSecond: listIntervalSecond,
	}
	c.run(ctx)
}
