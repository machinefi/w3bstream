package blockchain

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type contract struct {
	*monitor
	listInterval  time.Duration
	blockInterval uint64
}

func (t *contract) run(ctx context.Context) {
	ticker := time.NewTicker(t.listInterval)
	defer ticker.Stop()

	for range ticker.C {
		t.do(ctx)
	}
}

func (t *contract) do(ctx context.Context) {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.ContractLog{}

	_, l = l.Start(ctx, "contract.run")
	defer l.End()

	cs, err := m.List(d, builder.Or(
		m.ColBlockCurrent().Lt(m.ColBlockEnd()),
		m.ColBlockEnd().Eq(0),
	))
	if err != nil {
		l.Error(errors.Wrap(err, "list contractlog db failed"))
		return
	}
	for _, c := range cs {
		b := &models.Blockchain{RelBlockchain: models.RelBlockchain{ChainID: c.ChainID}}
		if err := b.FetchByChainID(d); err != nil {
			l.WithValues("chainID", c.ChainID).Error(errors.Wrap(err, "get chain info failed"))
			continue
		}
		toBlock, err := t.listChainAndSendEvent(ctx, &c, b.Address)
		if err != nil {
			l.Error(errors.Wrap(err, "list contractlog db failed"))
			continue
		}

		c.BlockCurrent = toBlock + 1
		if c.BlockEnd > 0 && c.BlockCurrent >= c.BlockEnd {
			c.Uniq = c.ContractLogID
		}
		if err := c.UpdateByID(d); err != nil {
			l.Error(errors.Wrap(err, "update contractlog db failed"))
		}
	}
}

func (t *contract) listChainAndSendEvent(ctx context.Context, c *models.ContractLog, address string) (uint64, error) {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "contract.listChainAndSendEvent")
	defer l.End()

	l = l.WithValues("type", "contract_log", "contract_log_id", c.ContractLogID)

	cli, err := ethclient.Dial(address)
	if err != nil {
		l.Error(err)
		return 0, err
	}

	from, to, err := t.getBlockRange(ctx, cli, c)
	if err != nil {
		l.Error(err)
		return 0, err
	}
	if from >= to {
		l.WithValues("from block", from, "to block", to).Debug("no new block")
		return to, nil
	}
	l.WithValues("from block", from, "to block", to).Debug("find new block")

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(from)),
		ToBlock:   big.NewInt(int64(to)),
		Addresses: []common.Address{
			common.HexToAddress(c.ContractAddress),
		},
		Topics: t.getTopic(c),
	}
	logs, err := cli.FilterLogs(context.Background(), query)
	if err != nil {
		l.Error(err)
		return 0, err
	}
	for _, log := range logs {
		data, err := log.MarshalJSON()
		if err != nil {
			return 0, err
		}
		if err := t.sendEvent(ctx, data, c.ProjectName, c.EventType); err != nil {
			return 0, err
		}
	}
	return to, nil
}

func (t *contract) getBlockRange(ctx context.Context, cli *ethclient.Client, c *models.ContractLog) (uint64, uint64, error) {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "contract.getBlockRange")
	defer l.End()

	currHeight, err := cli.BlockNumber(context.Background())
	if err != nil {
		l.Error(err)
		return 0, 0, err
	}
	from := c.BlockCurrent
	to := c.BlockCurrent + t.blockInterval
	if to > currHeight {
		to = currHeight
	}
	if c.BlockEnd > 0 && to > c.BlockEnd {
		to = c.BlockEnd
	}
	return from, to, nil
}

func (t *contract) getTopic(c *models.ContractLog) [][]common.Hash {
	res := make([][]common.Hash, 4)
	res[0] = t.parseTopic(c.Topic0)
	res[1] = t.parseTopic(c.Topic1)
	res[2] = t.parseTopic(c.Topic2)
	res[3] = t.parseTopic(c.Topic3)
	if len(res[3]) == 0 {
		res = res[:3]
		if len(res[2]) == 0 {
			res = res[:2]
			if len(res[1]) == 0 {
				res = res[:1]
				if len(res[0]) == 0 {
					res = res[:0]
				}
			}
		}
	}
	return res
}

func (t *contract) parseTopic(ts string) []common.Hash {
	res := make([]common.Hash, 0)
	if ts == "" {
		return res
	}
	ss := strings.Split(ts, ",")
	for _, s := range ss {
		h := common.HexToHash(s)
		res = append(res, h)
	}
	return res
}
