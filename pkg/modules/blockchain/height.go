package blockchain

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type height struct {
	*monitor
	interval time.Duration
}

func (h *height) run(ctx context.Context) {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	for range ticker.C {
		h.do(ctx)
	}
}

func (h *height) do(ctx context.Context) {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.ChainHeight{}

	_, l = l.Start(ctx, "height.run")
	defer l.End()

	cs, err := m.List(d, m.ColFinished().Eq(datatypes.FALSE))
	if err != nil {
		l.Error(errors.Wrap(err, "list chain height db failed"))
		return
	}
	for _, c := range cs {
		b := &models.Blockchain{RelBlockchain: models.RelBlockchain{ChainID: c.ChainID}}
		if err := b.FetchByChainID(d); err != nil {
			l.WithValues("chainID", c.ChainID).Error(errors.Wrap(err, "get chain info failed"))
			continue
		}
		res, err := h.checkHeightAndSendEvent(ctx, &c, b.Address)
		if err != nil {
			l.Error(errors.Wrap(err, "check chain height and send event failed"))
			continue
		}
		if res {
			c.Finished = datatypes.TRUE
			c.Uniq = c.ChainHeightID
			if err := c.UpdateByID(d); err != nil {
				l.Error(errors.Wrap(err, "update chain height db failed"))
			}
		}
	}
}

func (h *height) checkHeightAndSendEvent(ctx context.Context, c *models.ChainHeight, address string) (bool, error) {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "height.checkHeightAndSendEvent")
	defer l.End()

	l = l.WithValues("type", "chain_height", "chain_height_id", c.ChainHeightID)

	client, err := ethclient.Dial(address)
	if err != nil {
		l.Error(err)
		return false, err
	}
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		l.Error(err)
		return false, err
	}
	if headerNumber := header.Number.Uint64(); headerNumber < c.Height {
		l.WithValues("headerNumber", headerNumber, "chainHeight", c.Height).Debug("did not arrive")
		return false, nil
	}
	data, err := header.MarshalJSON()
	if err != nil {
		l.Error(err)
		return false, err
	}
	if err := h.sendEvent(ctx, data, c.ProjectName, c.EventType); err != nil {
		l.Error(err)
		return false, err
	}
	return true, nil
}
