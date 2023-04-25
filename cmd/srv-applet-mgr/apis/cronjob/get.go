package cronjob

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/cronjob"
	"github.com/machinefi/w3bstream/pkg/types"
)

type ListCronJob struct {
	httpx.MethodGet
	ProjectID       types.SFID `in:"path" name:"projectID"`
	cronjob.ListReq `in:"body"`
}

func (r *CreateCronJob) ListCronJob() string { return "/:projectID" }

func (r *ListCronJob) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextBySFID(ctx, r.ProjectID)
	if err != nil {
		return nil, err
	}
	r.ListReq.ProjectID = types.MustProjectFromContext(ctx).ProjectID
	return cronjob.List(ctx, &r.ListReq)
}
