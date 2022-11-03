package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

// ChainTx database model chain tx
// @def primary                   ID
// @def unique_index UI_chain_tx_id   ChainTxID
// @def unique_index UI_chain_tx_uniq ProjectName EventType ChainID TxAddress Uniq
//
//go:generate toolkit gen model ChainTx --database MonitorDB
type ChainTx struct {
	datatypes.PrimaryID
	RelChainTx
	ChainTxData
	datatypes.OperationTimes
}

type RelChainTx struct {
	ChainTxID types.SFID `db:"f_chaintx_id" json:"chaintxID"`
}

type ChainTxData struct {
	ProjectName string         `db:"f_project_name"                 json:"projectName"`
	Finished    datatypes.Bool `db:"f_finished,default='2'"         json:"-"`
	Uniq        types.SFID     `db:"f_uniq,default='0'"             json:"-"`
	ChainTxInfo
}

type ChainTxInfo struct {
	EventType string `db:"f_event_type"                   json:"eventType,omitempty"`
	ChainID   uint64 `db:"f_chain_id"                     json:"chainID"`
	TxAddress string `db:"f_tx_address"                   json:"txAddress"`
}
