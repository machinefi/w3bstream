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
	"github.com/iotexproject/Bumblebee/kit/sqlx"
	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/types"
)

const (
	listInterval  = 3 * time.Second
	blockInterval = 1000
)

func ListenContractLog(ctx context.Context) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.Blockchain{}
	ticker := time.NewTicker(listInterval)
	defer ticker.Stop()

	for range ticker.C {
		bs, err := m.List(d, nil)
		if err != nil {
			l.WithValues("info", "list blockchain db failed").Error(err)
			continue
		}
		for _, b := range bs {
			if err := listBlockChain(&b, d); err != nil {
				l.WithValues("info", "list blockchain failed").Error(err)
				continue
			}
		}
	}
}

func listBlockChain(bc *models.Blockchain, d sqlx.DBExecutor) error {
	client, err := ethclient.Dial(bc.BlockchainAddress)
	if err != nil {
		return err
	}

	contractAddress := common.HexToAddress(bc.ContractAddress)
	toBlock := bc.BlockCurrent + blockInterval
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(bc.BlockCurrent)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: []common.Address{
			contractAddress,
		},
	}
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://localhost:8888/srv-applet-mgr/v0/event/%s/%s/%s", bc.ProjectID, bc.AppletID, bc.Handler)

	for _, vLog := range logs {
		data, err := json.Marshal(vLog)
		if err != nil {
			return err
		}
		if err := sendEvent(data, url); err != nil {
			return err
		}
	}
	bc.BlockCurrent = toBlock
	return bc.UpdateByID(d)
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
