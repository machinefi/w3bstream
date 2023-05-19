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
	CronJobID types.SFID `in:"path" name:"cronJobID"`
}

func (r *RemoveCronJob) Path() string { return "/data/:cronJobID" }

func (r *RemoveCronJob) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithCronJobBySFID(ctx, r.CronJobID)
	if err != nil {
		return nil, err
	}
	return nil, cronjob.RemoveBySFID(ctx, r.CronJobID)
}
