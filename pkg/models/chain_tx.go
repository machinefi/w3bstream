package models

import (
	"github.com/machinefi/Bumblebee/base/types"
	"github.com/machinefi/Bumblebee/kit/sqlx/datatypes"
)

// Chaintx database model chaintx
// @def primary                   ID
// @def unique_index UI_chain_tx_id   ChaintxID
// @def unique_index UI_chain_tx_uniq ProjectName EventType ChainID TxAddress Uniq
//
//go:generate toolkit gen model Chaintx --database MonitorDB
type Chaintx struct {
	datatypes.PrimaryID
	RelChaintx
	ChaintxData
	datatypes.OperationTimes
}

type RelChaintx struct {
	ChaintxID types.SFID `db:"f_chaintx_id" json:"chaintxID"`
}

type ChaintxData struct {
	ProjectName string     `db:"f_project_name"                 json:"projectName"`
	Finished    bool       `db:"f_finished,default='false'"     json:"finished,omitempty"`
	Uniq        types.SFID `db:"f_uniq,default='0'"             json:"uniq,omitempty"`
	ChaintxInfo
}

type ChaintxInfo struct {
	EventType string `db:"f_event_type"                   json:"eventType,omitempty"`
	ChainID   uint64 `db:"f_chain_id"                     json:"chainID"`
	TxAddress string `db:"f_tx_address"                   json:"txAddress"`
}
