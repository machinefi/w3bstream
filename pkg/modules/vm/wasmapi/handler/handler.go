package handler

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Handler struct {
	mgrDB     sqlx.DBExecutor
	ethCli    *types.ETHClientConfig
	chainConf *types.ChainConfig
}

func New(mgrDB sqlx.DBExecutor, ethCli *types.ETHClientConfig, chainConf *types.ChainConfig) *Handler {
	return &Handler{
		mgrDB:     mgrDB,
		ethCli:    ethCli,
		chainConf: chainConf,
	}
}
