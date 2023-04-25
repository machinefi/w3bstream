package cronjob

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/cronjob"
	"github.com/machinefi/w3bstream/pkg/types"
)

type RemoveCronJob struct {
	httpx.MethodDelete
	ProjectID types.SFID `in:"path" name:"projectID"`
	CronJobID types.SFID `in:"path" name:"cronJobID"`
}

func (r *RemoveCronJob) Path() string { return "/:projectID/:cronJobID" }

func (r *RemoveCronJob) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	_, err := ca.WithProjectContextBySFID(ctx, r.ProjectID)
	if err != nil {
		return nil, err
	}
	return nil, cronjob.RemoveCronJob(ctx, r.ProjectID, r.CronJobID)
}
