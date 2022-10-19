package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/Bumblebee/conf/log"
	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"

	"github.com/iotexproject/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/iotexproject/w3bstream/pkg/enums"
	"github.com/iotexproject/w3bstream/pkg/errors/status"
	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/types"
)

const (
	listInterval  = 3 * time.Second
	blockInterval = 1000
)

func InitChainDB(ctx context.Context) error {
	d := types.MustDBExecutorFromContext(ctx)

	m := &models.Blockchain{
		RelBlockchain:  models.RelBlockchain{ChainID: 4690},
		BlockchainInfo: models.BlockchainInfo{Address: "https://babel-api.testnet.iotex.io"},
	}

	results := make([]models.Account, 0)
	err := d.QueryAndScan(builder.Select(nil).
		From(
			d.T(m),
			builder.Where(
				builder.And(
					m.ColChainID().Eq(4690),
				),
			),
		), &results)
	if err != nil {
		return status.CheckDatabaseError(err, "FetchChain")
	}
	if len(results) > 0 {
		return nil
	}
	return m.Create(d)
}

func ListenContractlog(ctx context.Context) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.Contractlog{}
	ticker := time.NewTicker(listInterval)
	defer ticker.Stop()

	for range ticker.C {
		cs, err := m.List(d, m.ColBlockCurrent().Lt(m.ColBlockEnd()))
		if err != nil {
			l.WithValues("info", "list contractlog db failed").Error(err)
			continue
		}
		for _, c := range cs {
			b := &models.Blockchain{RelBlockchain: models.RelBlockchain{ChainID: c.ChainID}}
			if err := b.FetchByChainID(d); err != nil {
				l.WithValues("info", "get chain info failed", "chainID", c.ChainID).Error(err)
				continue
			}
			toBlock, err := listChainAndSendEvent(l, &c, b.Address)
			if err != nil {
				l.WithValues("info", "list contractlog db failed").Error(err)
				continue
			}

			c.BlockCurrent = toBlock
			if err := c.UpdateByID(d); err != nil {
				l.WithValues("info", "update contractlog db failed").Error(err)
				continue
			}
		}
	}
}

func listChainAndSendEvent(l log.Logger, c *models.Contractlog, address string) (uint64, error) {
	cli, err := ethclient.Dial(address)
	if err != nil {
		return 0, err
	}
	ctx := context.Background()
	from, to, err := getBlockRange(ctx, cli, c)
	if err != nil {
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
		Topics: getTopic(c),
	}
	logs, err := cli.FilterLogs(ctx, query)
	if err != nil {
		return 0, err
	}
	url := fmt.Sprintf("http://localhost:8888/srv-applet-mgr/v0/event/%s", c.ProjectName)

	for _, vLog := range logs {
		data, err := json.Marshal(vLog)
		if err != nil {
			return 0, err
		}
		if err := sendEvent(data, url, c.EventType); err != nil {
			return 0, err
		}
	}
	return to, nil
}

func getBlockRange(ctx context.Context, cli *ethclient.Client, c *models.Contractlog) (uint64, uint64, error) {
	currHeight, err := cli.BlockNumber(context.Background())
	if err != nil {
		return 0, 0, err
	}
	from := c.BlockCurrent
	to := c.BlockCurrent + blockInterval
	if to > currHeight {
		to = currHeight
	}
	if c.BlockEnd > 0 && to > c.BlockEnd {
		to = c.BlockEnd
	}
	return from, to, nil
}

func getTopic(c *models.Contractlog) [][]common.Hash {
	res := make([][]common.Hash, 4)
	res[0] = parseTopic(c.Topic0)
	res[1] = parseTopic(c.Topic1)
	res[2] = parseTopic(c.Topic2)
	res[3] = parseTopic(c.Topic3)
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

func parseTopic(ts string) []common.Hash {
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

func sendEvent(data []byte, url string, et enums.EventType) error {
	// TODO event type
	e := &eventpb.Event{
		Payload: string(data),
	}
	body, err := json.Marshal(e)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// TODO http code judge
	return nil
}
