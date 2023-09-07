package handler

import (
	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	optypes "github.com/machinefi/w3bstream/pkg/modules/operator/pool/types"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Handler struct {
	opPool    optypes.Pool
	mgrDB     sqlx.DBExecutor
	chainConf *types.ChainConfig
	sfid      confid.SFIDGenerator
}

func New(mgrDB sqlx.DBExecutor, chainConf *types.ChainConfig, opPool optypes.Pool, sfid confid.SFIDGenerator) *Handler {
	return &Handler{
		opPool:    opPool,
		mgrDB:     mgrDB,
		chainConf: chainConf,
		sfid:      sfid,
	}
}
