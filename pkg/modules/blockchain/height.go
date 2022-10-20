package blockchain

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/Bumblebee/conf/log"

	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type height struct {
	*monitor
	interval time.Duration
}

func (h *height) run(ctx context.Context) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.ChainHeight{}
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	for range ticker.C {
		cs, err := m.List(d, m.ColFinished().Eq(false))
		if err != nil {
			l.WithValues("info", "list chain height db failed").Error(err)
			continue
		}
		for _, c := range cs {
			b := &models.Blockchain{RelBlockchain: models.RelBlockchain{ChainID: c.ChainID}}
			if err := b.FetchByChainID(d); err != nil {
				l.WithValues("info", "get chain info failed", "chainID", c.ChainID).Error(err)
				continue
			}
			res, err := h.checkHeightAndSendEvent(l, &c, b.Address)
			if err != nil {
				l.WithValues("info", "check chain height and send event failed").Error(err)
				continue
			}
			if res {
				c.Finished = true
				if err := c.UpdateByID(d); err != nil {
					l.WithValues("info", "update chain height db failed").Error(err)
				}
			}
		}
	}
}

func (h *height) checkHeightAndSendEvent(l log.Logger, c *models.ChainHeight, address string) (bool, error) {
	client, err := ethclient.Dial(address)
	if err != nil {
		return false, err
	}
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return false, err
	}
	if header.Number.Uint64() < c.Height {
		return false, nil
	}
	data, err := header.MarshalJSON()
	if err != nil {
		return false, err
	}
	if err := h.sendEvent(data, c.ProjectName, c.EventType); err != nil {
		return false, err
	}
	return true, nil
}
