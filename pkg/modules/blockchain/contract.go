package blockchain

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type contract struct {
	*monitor
	listInterval  time.Duration
	blockInterval uint64
}

// a group of models.ContractLog, will list the chain and update current height togither
type listChainGroup struct {
	toBlock uint64
	cs      []*models.ContractLog
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

	cs, err := m.List(d, builder.And(
		builder.Or(
			m.ColBlockCurrent().Lt(m.ColBlockEnd()),
			m.ColBlockEnd().Eq(0),
		),
		m.ColPaused().Eq(datatypes.FALSE),
	))
	if err != nil {
		l.Error(errors.Wrap(err, "list contractlog db failed"))
		return
	}

	gs, err := t.getListChainGroups(ctx, cs)
	if err != nil {
		l.Error(errors.Wrap(err, "get lister units failed"))
		return
	}

	for _, g := range gs {
		toBlock, err := t.listChainAndSendEvent(ctx, g)
		if err != nil {
			l.Error(errors.Wrap(err, "list chain and send event failed"))
			continue
		}

		if err := sqlx.NewTasks(d).With(
			func(d sqlx.DBExecutor) error {
				for _, c := range g.cs {
					c.BlockCurrent = toBlock + 1
					if c.BlockEnd > 0 && c.BlockCurrent >= c.BlockEnd {
						c.Uniq = c.ContractLogID
					}
					if err := c.UpdateByID(d); err != nil {
						return err
					}
				}
				return nil
			},
		).Do(); err != nil {
			l.Error(errors.Wrap(err, "update contractlog db failed"))
		}
	}
}

func (t *contract) getListChainGroups(ctx context.Context, cs []models.ContractLog) ([]*listChainGroup, error) {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "contract.getListChainGroups")
	defer l.End()

	us := t.groupContractLog(cs)
	t.pruneListChainGroups(us)
	if err := t.setToBlock(ctx, us); err != nil {
		return nil, err
	}
	return us, nil
}

// projectName + chainID -> contractLog list
func (t *contract) groupContractLog(cs []models.ContractLog) []*listChainGroup {
	groups := make(map[string][]*models.ContractLog)

	for i := range cs {
		key := fmt.Sprintf("%s_%d", cs[i].ProjectName, cs[i].ChainID)
		groups[key] = append(groups[key], &cs[i])
	}

	ret := []*listChainGroup{}
	for _, cs := range groups {
		ret = append(ret, &listChainGroup{
			cs: cs,
		})
	}
	return ret
}

func (t *contract) pruneListChainGroups(gs []*listChainGroup) {
	for _, g := range gs {
		sort.SliceStable(g.cs, func(i, j int) bool {
			return g.cs[i].BlockCurrent < g.cs[j].BlockCurrent
		})

		if g.cs[0].BlockCurrent == g.cs[len(g.cs)-1].BlockCurrent {
			continue
		}
		for i := range g.cs {
			if i == 0 {
				continue
			}
			if g.cs[i].BlockCurrent != g.cs[i-1].BlockCurrent {
				g.toBlock = g.cs[i].BlockCurrent - 1
				g.cs = g.cs[:i]
				break
			}
		}
	}
}

func (t *contract) setToBlock(ctx context.Context, gs []*listChainGroup) error {
	l := types.MustLoggerFromContext(ctx)
	ethcli := types.MustETHClientConfigFromContext(ctx)

	_, l = l.Start(ctx, "contract.setToBlock")
	defer l.End()

	for _, g := range gs {
		c := g.cs[0]

		chainAddress, ok := ethcli.Clients[uint32(c.ChainID)]
		if !ok {
			err := errors.New("blockchain not exist")
			l.WithValues("chainID", c.ChainID).Error(err)
			return err
		}

		cli, err := ethclient.Dial(chainAddress)
		if err != nil {
			l.WithValues("chainID", c.ChainID).Error(errors.Wrap(err, "dial eth address failed"))
			return err
		}
		currHeight, err := cli.BlockNumber(context.Background())
		if err != nil {
			l.Error(errors.Wrap(err, "get blockchain current height failed"))
			return err
		}

		to := c.BlockCurrent + t.blockInterval
		if to > currHeight {
			to = currHeight
		}
		for _, c := range g.cs {
			if c.BlockEnd > 0 && to > c.BlockEnd {
				to = c.BlockEnd
			}
		}
		if g.toBlock == 0 {
			g.toBlock = to
		}
		if g.toBlock > to {
			g.toBlock = to
		}
	}
	return nil
}

func (t *contract) listChainAndSendEvent(ctx context.Context, g *listChainGroup) (uint64, error) {
	l := types.MustLoggerFromContext(ctx)
	ethcli := types.MustETHClientConfigFromContext(ctx)

	_, l = l.Start(ctx, "contract.listChainAndSendEvent")
	defer l.End()

	c := g.cs[0]

	l = l.WithValues("chainID", c.ChainID, "projectName", c.ProjectName)

	chainAddress, ok := ethcli.Clients[uint32(c.ChainID)]
	if !ok {
		err := errors.New("blockchain not exist")
		l.Error(err)
		return 0, err
	}

	cli, err := ethclient.Dial(chainAddress)
	if err != nil {
		l.Error(errors.Wrap(err, "dial eth address failed"))
		return 0, err
	}

	from, to := c.BlockCurrent, g.toBlock

	if from > to {
		l.WithValues("from block", from, "to block", to).Debug("no new block")
		return to, nil
	}
	l.WithValues("from block", from, "to block", to).Debug("find new block")

	as, mas := t.getAddresses(g.cs)
	ts, mts := t.getTopic(g.cs)
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(from)),
		ToBlock:   big.NewInt(int64(to)),
		Addresses: as,
		Topics:    ts,
	}
	logs, err := cli.FilterLogs(context.Background(), query)
	if err != nil {
		l.Error(errors.Wrap(err, "filter event logs failed"))
		return 0, err
	}
	for i := range logs {
		cs := t.getExpectedContractLogs(&logs[i], mas, mts)
		if len(cs) == 0 {
			err := errors.New("cannot find expected contract log")
			l.Error(err)
			return 0, err
		}

		data, err := logs[i].MarshalJSON()
		if err != nil {
			return 0, err
		}
		for _, c := range cs {
			if err := t.sendEvent(ctx, data, c.ProjectName, c.EventType); err != nil {
				return 0, err
			}
		}
	}
	return to, nil
}

func (t *contract) getExpectedContractLogs(log *ethtypes.Log, mas map[*models.ContractLog]common.Address, mts map[*models.ContractLog][]*common.Hash) []*models.ContractLog {
	res := []*models.ContractLog{}

	for c, addr := range mas {
		if bytes.Equal(addr.Bytes(), log.Address.Bytes()) {
			ts := mts[c]

			for i, contractLogTopic := range ts {
				if contractLogTopic == nil {
					continue
				}
				if len(log.Topics) > i && bytes.Equal(log.Topics[i].Bytes(), contractLogTopic.Bytes()) {
					continue
				}
				goto Next
			}
			res = append(res, c)
		}
	Next:
	}
	return res
}

func (t *contract) getAddresses(cs []*models.ContractLog) ([]common.Address, map[*models.ContractLog]common.Address) {
	as := []common.Address{}
	mas := make(map[*models.ContractLog]common.Address)
	for _, c := range cs {
		a := common.HexToAddress(c.ContractAddress)
		as = append(as, a)
		mas[c] = a
	}
	return as, mas
}

func (t *contract) getTopic(cs []*models.ContractLog) ([][]common.Hash, map[*models.ContractLog][]*common.Hash) {
	res := make([][]common.Hash, 4)
	mres := make(map[*models.ContractLog][]*common.Hash)

	for _, c := range cs {
		h0 := t.parseTopic(c.Topic0)
		mres[c] = append(mres[c], h0)
		if h0 != nil {
			res[0] = append(res[0], *h0)
		}

		h1 := t.parseTopic(c.Topic1)
		mres[c] = append(mres[c], h1)
		if h1 != nil {
			res[1] = append(res[1], *h1)
		}

		h2 := t.parseTopic(c.Topic2)
		mres[c] = append(mres[c], h2)
		if h2 != nil {
			res[2] = append(res[2], *h2)
		}

		h3 := t.parseTopic(c.Topic3)
		mres[c] = append(mres[c], h3)
		if h3 != nil {
			res[3] = append(res[3], *h3)
		}

	}

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
	return res, mres
}

func (t *contract) parseTopic(tStr string) *common.Hash {
	if tStr == "" {
		return nil
	}
	h := common.HexToHash(tStr)
	return &h
}
