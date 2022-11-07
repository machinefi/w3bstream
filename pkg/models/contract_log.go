package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

// ContractLog database model contract log
// @def primary                   ID
// @def unique_index UI_contract_log_id   ContractLogID
// @def unique_index UI_contract_log_uniq ProjectName EventType ChainID ContractAddress Topic0 Topic1 Topic2 Topic3 Uniq
//
//go:generate toolkit gen model ContractLog --database MonitorDB
type ContractLog struct {
	datatypes.PrimaryID
	RelContractLog
	ContractLogData
	datatypes.OperationTimes
}

type RelContractLog struct {
	ContractLogID types.SFID `db:"f_contractlog_id" json:"contractlogID"`
}

type ContractLogData struct {
	ProjectName string     `db:"f_project_name"                 json:"projectName"`
	Uniq        types.SFID `db:"f_uniq,default='0'"             json:"-"`
	ContractLogInfo
}

type ContractLogInfo struct {
	EventType       string `db:"f_event_type"                   json:"eventType,omitempty"`
	ChainID         uint64 `db:"f_chain_id"                     json:"chainID"`
	ContractAddress string `db:"f_contract_address"             json:"contractAddress"`
	BlockStart      uint64 `db:"f_block_start"                  json:"blockStart"`
	BlockCurrent    uint64 `db:"f_block_current"                json:"blockCurrent,omitempty"`
	BlockEnd        uint64 `db:"f_block_end,default='0'"        json:"blockEnd,omitempty"`
	Topic0          string `db:"f_topic0,default=''"            json:"topic0,omitempty"`
	Topic1          string `db:"f_topic1,default=''"            json:"topic1,omitempty"`
	Topic2          string `db:"f_topic2,default=''"            json:"topic2,omitempty"`
	Topic3          string `db:"f_topic3,default=''"            json:"topic3,omitempty"`
}
