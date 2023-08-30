package handler

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	optypes "github.com/machinefi/w3bstream/pkg/modules/operator/pool/types"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Handler struct {
	opPool    optypes.Pool
	mgrDB     sqlx.DBExecutor
	chainConf *types.ChainConfig
}

func New(mgrDB sqlx.DBExecutor, chainConf *types.ChainConfig, opPool optypes.Pool) *Handler {
	return &Handler{
		opPool:    opPool,
		mgrDB:     mgrDB,
		chainConf: chainConf,
	}
}
