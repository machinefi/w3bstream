package configuration

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/types"
)

type EthClientRsp struct {
	Clients string `json:"clients"`
}

type EthClient struct {
	httpx.MethodGet
}

func (r *EthClient) Path() string { return "/eth_client" }

func (r *EthClient) Output(ctx context.Context) (interface{}, error) {
	ethcli := types.MustETHClientConfigFromContext(ctx)
	return &EthClientRsp{
		Clients: ethcli.Endpoints,
	}, nil
}

type ChainConfigResp struct {
	Chains string `json:"chains"`
}

type ChainConfig struct {
	httpx.MethodGet
}

func (r *ChainConfig) Path() string { return "/chain_config" }

func (r *ChainConfig) Output(ctx context.Context) (interface{}, error) {
	c := types.MustChainConfigFromContext(ctx)
	return &ChainConfigResp{
		Chains: c.Configs,
	}, nil
}
