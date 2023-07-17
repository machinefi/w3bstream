package handler

import (
	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Handler struct {
	l      log.Logger
	mgrDB  sqlx.DBExecutor
	ethCli *types.ETHClientConfig
}

func New(mgrDB sqlx.DBExecutor, l log.Logger, ethCli *types.ETHClientConfig) *Handler {
	return &Handler{
		l:      l,
		mgrDB:  mgrDB,
		ethCli: ethCli,
	}
}
