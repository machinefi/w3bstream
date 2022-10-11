package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://babel-api.testnet.iotex.io")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0xb93Fc2a4729C9EF8Bd202Bc70A04c19654D78d57")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(16737070),
		ToBlock:   big.NewInt(16737073),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	url := "http://localhost:8888/srv-applet-mgr/v0/event/{project_id}/{applet_id}/start"

	for _, vLog := range logs {
		info, err := json.Marshal(vLog)
		if err != nil {
			panic(err)
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(info))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("publisher", "test publisher")

		cli := &http.Client{}
		resp, err := cli.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	}
}
