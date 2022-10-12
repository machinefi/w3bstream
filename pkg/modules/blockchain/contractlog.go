package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/types"
)

const (
	listInterval  = 3 * time.Second
	blockInterval = 1000
)

func ListenContractlog(ctx context.Context) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.Contractlog{}
	ticker := time.NewTicker(listInterval)
	defer ticker.Stop()

	for range ticker.C {
		cs, err := m.List(d, nil)
		if err != nil {
			l.WithValues("info", "list contractlog db failed").Error(err)
			continue
		}
		for _, c := range cs {
			toBlock, err := listChainAndSendEvent(ctx, &c)
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

func listChainAndSendEvent(ctx context.Context, c *models.Contractlog) (uint64, error) {
	l := types.MustLoggerFromContext(ctx)

	var address string // TODO howto get address by chainID
	client, err := ethclient.Dial(address)
	if err != nil {
		return 0, err
	}

	from, to, err := getBlockRange(client, c)
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
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return 0, err
	}
	url := fmt.Sprintf("http://localhost:8888/srv-applet-mgr/v0/event/%s", c.ProjectName)

	for _, vLog := range logs {
		data, err := json.Marshal(vLog)
		if err != nil {
			return 0, err
		}
		if err := sendEvent(data, url); err != nil {
			return 0, err
		}
	}
	return to, nil
}

func getBlockRange(cli *ethclient.Client, c *models.Contractlog) (uint64, uint64, error) {
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
	t1 := make([]common.Hash, 0)
	if c.Topic1 != "" {
		h1 := common.HexToHash(c.Topic1)
		t1 = append(t1, h1)
	}
	t2 := make([]common.Hash, 0)
	if c.Topic2 != "" {
		h2 := common.HexToHash(c.Topic2)
		t2 = append(t2, h2)
	}
	t3 := make([]common.Hash, 0)
	if c.Topic3 != "" {
		h3 := common.HexToHash(c.Topic3)
		t3 = append(t3, h3)
	}
	return append(make([][]common.Hash, 0, 3), t1, t2, t3)
}

func sendEvent(data []byte, url string) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("publisher", "test publisher") // TODO set publisher

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// TODO http code judge
	return nil
}
