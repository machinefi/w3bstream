package cronjob

import (
	"context"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateCronJobReq = models.CronJobInfo

func CreateCronJob(ctx context.Context, projectID types.SFID, r *CreateCronJobReq) (*models.CronJob, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "CreateCronJob")
	defer l.End()

	if _, err := cron.ParseStandard(r.CronExpressions); err != nil {
		l.WithValues("cronExpressions", r.CronExpressions).Error(errors.Wrap(err, "cron expressions illegal"))
		return nil, status.BadRequest.StatusErr().WithDesc("cron expressions illegal")
	}

	n := *r
	n.EventType = getEventType(n.EventType)
	m := &models.CronJob{
		RelCronJob:  models.RelCronJob{CronJobID: idg.MustGenSFID()},
		RelProject:  models.RelProject{ProjectID: projectID},
		CronJobInfo: n,
	}
	if err := m.Create(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "CreateCronJob")
	}
	return m, nil
}

func RemoveCronJob(ctx context.Context, projectID, cronJobID types.SFID) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "RemoveCronJob")
	defer l.End()

	l = l.WithValues("cronJobID", cronJobID)

	m := &models.CronJob{RelCronJob: models.RelCronJob{CronJobID: cronJobID}}
	if err := m.FetchByCronJobID(d); err != nil {
		return status.CheckDatabaseError(err, "FetchByCronJobID")
	}
	if err := checkProjectID(m.ProjectID, projectID, l); err != nil {
		return err
	}
	if err := m.DeleteByCronJobID(d); err != nil {
		return status.CheckDatabaseError(err, "DeleteByCronJobID")
	}
	return nil
}

func checkProjectID(want, curr types.SFID, l log.Logger) error {
	if want != curr {
		l.Error(errors.New("cron job project mismatch"))
		return status.BadRequest.StatusErr().WithDesc("cron job project mismatch")
	}
	return nil
}

func getEventType(eventType string) string {
	if eventType == "" {
		return enums.EVENTTYPEDEFAULT
	}
	return eventType
}
