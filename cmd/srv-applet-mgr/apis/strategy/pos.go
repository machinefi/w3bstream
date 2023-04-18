package strategy

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateStrategy struct {
	httpx.MethodPost
	strategy.BatchCreateReq `in:"body"`
}

func (r *CreateStrategy) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, middleware.ProjectProviderFromContext(ctx))
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)
	ret := make([]models.Strategy, 0, len(r.Strategies))
	ids := confid.MustSFIDGeneratorFromContext(ctx).MustGenSFIDs(len(r.Strategies))

	for i := range r.Strategies {
		ret = append(ret, models.Strategy{
			RelStrategy:  models.RelStrategy{StrategyID: ids[i]},
			RelProject:   models.RelProject{ProjectID: prj.ProjectID},
			RelApplet:    r.Strategies[i].RelApplet,
			StrategyInfo: r.Strategies[i].StrategyInfo,
		})
	}

	if err = strategy.BatchCreate(ctx, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
