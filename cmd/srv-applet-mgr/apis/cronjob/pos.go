package cronjob

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/cronjob"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateCronJob struct {
	httpx.MethodPost
	ProjectID                types.SFID `in:"path" name:"projectID"`
	cronjob.CreateCronJobReq `in:"body"`
}

func (r *CreateCronJob) Path() string { return "/:projectID" }

func (r *CreateCronJob) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	_, err := ca.WithProjectContextByID(ctx, r.ProjectID)
	if err != nil {
		return nil, err
	}
	return cronjob.CreateCronJob(ctx, r.ProjectID, &r.CreateCronJobReq)
}
