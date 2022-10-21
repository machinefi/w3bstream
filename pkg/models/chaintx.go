package models

import (
	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/sqlx/datatypes"

	"github.com/iotexproject/w3bstream/pkg/enums"
)

// Chaintx database model chaintx
// @def primary                   ID
//
//go:generate toolkit gen model Chaintx --database DB
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
	ProjectName string `db:"f_project_name"                 json:"projectName"`
	Finished    bool   `db:"f_finished,default='false'"     json:"finished,omitempty"`
	ChaintxInfo
}

type ChaintxInfo struct {
	EventType enums.EventType `db:"f_event_type"                   json:"eventType,omitempty"`
	ChainID   uint64          `db:"f_chain_id"                     json:"chainID"`
	TxAddress string          `db:"f_tx_address"                   json:"txAddress"`
}
