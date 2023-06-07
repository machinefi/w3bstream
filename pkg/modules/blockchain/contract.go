package blockchain

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type contract struct {
	*monitor
	listInterval  time.Duration
	blockInterval uint64
}

type listerUnit struct {
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

	cs, err := m.List(d, builder.Or(
		m.ColBlockCurrent().Lt(m.ColBlockEnd()),
		m.ColBlockEnd().Eq(0),
	))
	if err != nil {
		l.Error(errors.Wrap(err, "list contractlog db failed"))
		return
	}

	us, err := t.getListerUnits(ctx, cs)
	if err != nil {
		l.Error(errors.Wrap(err, "get lister units failed"))
		return
	}

	for _, u := range us {
		toBlock, err := t.listChainAndSendEvent(ctx, u)
		if err != nil {
			l.Error(errors.Wrap(err, "list chain and send event failed"))
			continue
		}

		if err := sqlx.NewTasks(d).With(
			func(d sqlx.DBExecutor) error {
				for _, c := range u.cs {
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

func (t *contract) getListerUnits(ctx context.Context, cs []models.ContractLog) ([]*listerUnit, error) {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "contract.getListerUnits")
	defer l.End()

	us := t.classifyContractLog(cs)
	t.pruneListerUnits(us)
	if err := t.setToBlock(ctx, us); err != nil {
		return nil, err
	}
	return us, nil
}

// projectName + chainID -> contractLog list
func (t *contract) classifyContractLog(cs []models.ContractLog) []*listerUnit {
	class := make(map[string][]*models.ContractLog)

	for i := range cs {
		key := fmt.Sprintf("%s_%d", cs[i].ProjectName, cs[i].ChainID)
		class[key] = append(class[key], &cs[i])
	}

	ret := []*listerUnit{}
	for _, cs := range class {
		ret = append(ret, &listerUnit{
			cs: cs,
		})
	}
	return ret
}

func (t *contract) pruneListerUnits(us []*listerUnit) {
	for _, u := range us {
		sort.SliceStable(u.cs, func(i, j int) bool {
			return u.cs[i].BlockCurrent < u.cs[j].BlockCurrent
		})

		if u.cs[0].BlockCurrent == u.cs[len(u.cs)-1].BlockCurrent {
			continue
		}
		for i := range u.cs {
			if i == 0 {
				continue
			}
			if u.cs[i].BlockCurrent != u.cs[i-1].BlockCurrent {
				u.toBlock = u.cs[i].BlockCurrent - 1
				u.cs = u.cs[:i]
				break
			}
		}
	}
}

func (t *contract) setToBlock(ctx context.Context, us []*listerUnit) error {
	l := types.MustLoggerFromContext(ctx)
	ethcli := types.MustETHClientConfigFromContext(ctx)

	_, l = l.Start(ctx, "contract.setToBlock")
	defer l.End()

	for _, u := range us {
		c := u.cs[0]

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
		for _, c := range u.cs {
			if c.BlockEnd > 0 && to > c.BlockEnd {
				to = c.BlockEnd
			}
		}
		if u.toBlock == 0 {
			u.toBlock = to
		}
		if u.toBlock > to {
			u.toBlock = to
		}
	}
	return nil
}

func (t *contract) listChainAndSendEvent(ctx context.Context, u *listerUnit) (uint64, error) {
	l := types.MustLoggerFromContext(ctx)
	ethcli := types.MustETHClientConfigFromContext(ctx)

	_, l = l.Start(ctx, "contract.listChainAndSendEvent")
	defer l.End()

	c := u.cs[0]

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

	from, to := c.BlockCurrent, u.toBlock

	if from > to {
		l.WithValues("from block", from, "to block", to).Debug("no new block")
		return to, nil
	}
	l.WithValues("from block", from, "to block", to).Debug("find new block")

	as, mas := t.getAddresses(u.cs)
	ts, mts := t.getTopic(u.cs)
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
		c, err := t.getExpectedContractLog(&logs[i], mas, mts)
		if err != nil {
			l.Error(err)
			return 0, err
		}

		data, err := logs[i].MarshalJSON()
		if err != nil {
			return 0, err
		}
		if err := t.sendEvent(ctx, data, c.ProjectName, c.EventType); err != nil {
			return 0, err
		}
	}
	return to, nil
}

func (t *contract) getExpectedContractLog(log *ethtypes.Log, mas map[*models.ContractLog]common.Address, mts map[*models.ContractLog][]common.Hash) (*models.ContractLog, error) {
	logTopics := make(map[string]bool)
	for _, l := range log.Topics {
		logTopics[l.String()] = true
	}

	for c, as := range mas {
		if bytes.Equal(as.Bytes(), log.Address.Bytes()) {
			ts := mts[c]
			for _, t := range ts {
				if _, ok := logTopics[t.String()]; !ok {
					goto Next
				}
			}
			return c, nil
		}
	Next:
	}
	return nil, errors.New("cannot find expected contract log")
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

func (t *contract) getTopic(cs []*models.ContractLog) ([][]common.Hash, map[*models.ContractLog][]common.Hash) {
	res := make([][]common.Hash, 4)
	mres := make(map[*models.ContractLog][]common.Hash)

	for _, c := range cs {
		t0 := t.parseTopic(c.Topic0)
		t1 := t.parseTopic(c.Topic1)
		t2 := t.parseTopic(c.Topic2)
		t3 := t.parseTopic(c.Topic3)

		res[0] = append(res[0], t0...)
		res[1] = append(res[1], t1...)
		res[2] = append(res[2], t2...)
		res[3] = append(res[3], t3...)

		mres[c] = append(mres[c], t0...)
		mres[c] = append(mres[c], t1...)
		mres[c] = append(mres[c], t2...)
		mres[c] = append(mres[c], t3...)
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
