package ratelimit

import (
	"context"
	"time"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/job"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateTrafficRateLimitReq struct {
	models.RateLimitInfo
}

func CreateRateLimit(ctx context.Context, r *CreateTrafficRateLimitReq) (*models.TrafficRateLimit, error) {
	//d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)
	project := types.MustProjectFromContext(ctx)

	_, l = l.Start(ctx, "CreateTrafficRateLimit")
	defer l.End()

	m := &models.TrafficRateLimit{
		RelRateLimit: models.RelRateLimit{RateLimitID: idg.MustGenSFID()},
		RelProject:   models.RelProject{ProjectID: project.ProjectID},
		RateLimitInfo: models.RateLimitInfo{
			Count:    r.Count,
			Duration: r.Duration,
			ApiType:  r.ApiType,
		},
	}
	//if err := m.Create(d); err != nil {
	//	l.Error(err)
	//	return nil, err
	//}

	t := job.NewTrafficTask(*m)
	job.Dispatch(ctx, t)
	time.Sleep(3 * time.Second)

	return m, nil
}

func UpdateRateLimit() (*models.TrafficRateLimit, error) {
	//
	return nil, nil
}
