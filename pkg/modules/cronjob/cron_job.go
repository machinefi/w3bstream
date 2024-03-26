package cronjob

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
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
	ctx, l := logr.Start(ctx, "cronjob.run")
	defer l.End()

	s := gocron.NewScheduler(time.UTC)
	s.TagsUnique()

	if _, err := s.Every(t.listIntervalSecond).Seconds().Do(t.do, ctx, s); err != nil {
		l.Error(errors.Wrap(err, "create cronjob main loop failed"))
	}
	s.StartAsync()
}

func (t *cronJob) do(ctx context.Context, s *gocron.Scheduler) {
	ctx, l := logger.NewSpanContext(ctx, "cronjob.do")
	defer l.End()

	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.CronJob{}

	cs, err := m.List(d, nil)
	if err != nil {
		l.Error(errors.Wrap(err, "list cronjob db failed"))
		return
	}

	t.tidyCronJobs(ctx, s, cs)

	for _, c := range cs {
		if _, err := s.Cron(c.CronExpressions).Tag(c.CronJobID.String()).Do(t.sendEvent, ctx, c); err != nil {
			if !strings.Contains(err.Error(), "non-unique tag") {
				l.WithValues("cronJobID", c.CronJobID).Error(errors.Wrap(err, "create new cron job failed"))
			}
		}
	}
}

func (t *cronJob) tidyCronJobs(ctx context.Context, s *gocron.Scheduler, cs []models.CronJob) {
	ctx, l := logr.Start(ctx, "cronjob.tidyCronJobs")
	defer l.End()

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
	ctx, l := logr.Start(ctx, "cronjob.sendEvent", "cronJobID", c.CronJobID)
	defer l.End()

	d := types.MustMgrDBExecutorFromContext(ctx)

	m := &models.Project{RelProject: models.RelProject{ProjectID: c.ProjectID}}
	if err := m.FetchByProjectID(d); err != nil {
		l.Error(errors.Wrap(err, "get project failed"))
		return
	}
	payload, err := json.Marshal(struct {
		ID        types.SFID
		Timestamp time.Time
	}{
		c.CronJobID,
		time.Now(),
	})
	if err != nil {
		l.Error(errors.Wrap(err, "failed to marshal payload"))
		return
	}

	ctx = types.WithProject(ctx, m)
	if _, err = event.HandleEvent(ctx, c.EventType, payload); err != nil {
		l.Error(errors.Wrap(err, "send event failed"))
	}
}

func Run(ctx context.Context) {
	c := &cronJob{
		listIntervalSecond: listIntervalSecond,
	}
	c.run(ctx)
}
