package blockchain

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type tx struct {
	*monitor
	interval time.Duration
}

func (t *tx) run(ctx context.Context) {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	for range ticker.C {
		t.do(ctx)
	}
}

func (t *tx) do(ctx context.Context) {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.Chaintx{}

	_, l = l.Start(ctx, "tx.run")
	defer l.End()

	cs, err := m.List(d, m.ColFinished().Eq(false))
	if err != nil {
		l.Error(errors.Wrap(err, "list chain tx db failed"))
		return
	}
	for _, c := range cs {
		b := &models.Blockchain{RelBlockchain: models.RelBlockchain{ChainID: c.ChainID}}
		if err := b.FetchByChainID(d); err != nil {
			l.WithValues("chainID", c.ChainID).Error(errors.Wrap(err, "get chain info failed"))
			return
		}
		res, err := t.checkTxAndSendEvent(ctx, &c, b.Address)
		if err != nil {
			l.Error(errors.Wrap(err, "check chain tx and send event failed"))
			return
		}
		if res {
			c.Finished = true
			if err := c.UpdateByID(d); err != nil {
				l.Error(errors.Wrap(err, "update chain tx db failed"))
			}
		}
	}

}

func (t *tx) checkTxAndSendEvent(ctx context.Context, c *models.Chaintx, address string) (bool, error) {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "tx.checkTxAndSendEvent")
	defer l.End()

	l = l.WithValues("type", "chain_tx", "chain_tx_id", c.ChaintxID)

	client, err := ethclient.Dial(address)
	if err != nil {
		l.Error(err)
		return false, err
	}
	tx, p, err := client.TransactionByHash(context.Background(), common.HexToHash(c.TxAddress))
	if err != nil {
		l.Error(err)
		return false, err
	}
	if p {
		l.WithValues("tx_hash", c.TxAddress).Debug("transaction pending")
		return false, nil
	}
	data, err := tx.MarshalJSON()
	if err != nil {
		l.Error(err)
		return false, err
	}
	if err := t.sendEvent(ctx, data, c.ProjectName, c.EventType); err != nil {
		l.Error(err)
		return false, err
	}
	return true, nil
}
