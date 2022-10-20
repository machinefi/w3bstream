package blockchain

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/Bumblebee/conf/log"

	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type tx struct {
	*monitor
	interval time.Duration
}

func (t *tx) run(ctx context.Context) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.Chaintx{}
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	for range ticker.C {
		cs, err := m.List(d, m.ColFinished().Eq(false))
		if err != nil {
			l.WithValues("info", "list chain tx db failed").Error(err)
			continue
		}
		for _, c := range cs {
			b := &models.Blockchain{RelBlockchain: models.RelBlockchain{ChainID: c.ChainID}}
			if err := b.FetchByChainID(d); err != nil {
				l.WithValues("info", "get chain info failed", "chainID", c.ChainID).Error(err)
				continue
			}
			res, err := t.checkTxAndSendEvent(l, &c, b.Address)
			if err != nil {
				l.WithValues("info", "check chain tx and send event failed").Error(err)
				continue
			}
			if res {
				c.Finished = true
				if err := c.UpdateByID(d); err != nil {
					l.WithValues("info", "update chain tx db failed").Error(err)
				}
			}
		}
	}
}

func (t *tx) checkTxAndSendEvent(l log.Logger, c *models.Chaintx, address string) (bool, error) {
	client, err := ethclient.Dial(address)
	if err != nil {
		return false, err
	}
	tx, p, err := client.TransactionByHash(context.Background(), common.HexToHash(c.TxAddress))
	if err != nil {
		return false, err
	}
	if p {
		return false, nil
	}
	data, err := tx.MarshalJSON()
	if err != nil {
		return false, err
	}
	if err := t.sendEvent(data, c.ProjectName, c.EventType); err != nil {
		return false, err
	}
	return true, nil
}
