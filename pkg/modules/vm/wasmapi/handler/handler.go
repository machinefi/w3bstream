package handler

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Handler struct {
	mgrDB     sqlx.DBExecutor
	chainConf *types.ChainConfig
}

func New(mgrDB sqlx.DBExecutor, chainConf *types.ChainConfig) *Handler {
	return &Handler{
		mgrDB:     mgrDB,
		chainConf: chainConf,
	}
}
