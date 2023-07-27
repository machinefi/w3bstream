package handler

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Handler struct {
	mgrDB  sqlx.DBExecutor
	ethCli *types.ETHClientConfig
}

func New(mgrDB sqlx.DBExecutor, ethCli *types.ETHClientConfig) *Handler {
	return &Handler{
		mgrDB:  mgrDB,
		ethCli: ethCli,
	}
}
