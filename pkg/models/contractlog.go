package models

import (
	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/sqlx/datatypes"

	"github.com/iotexproject/w3bstream/pkg/enums"
)

// Contractlog database model contractlog
// @def primary                   ID
//
//go:generate toolkit gen model Contractlog --database DB
type Contractlog struct {
	datatypes.PrimaryID
	RelContractlog
	ContractlogData
	datatypes.OperationTimes
}

type RelContractlog struct {
	ContractlogID types.SFID `db:"f_contractlog_id" json:"contractlogID"`
}

type ContractlogData struct {
	ProjectName string `db:"f_project_name"                 json:"projectName"`
	ContractlogInfo
}

type ContractlogInfo struct {
	EventType       enums.EventType `db:"f_event_type"                   json:"eventType,omitempty"`
	ChainID         uint64          `db:"f_chain_id"                     json:"chainID"`
	ContractAddress string          `db:"f_contract_address"             json:"contractAddress"`
	BlockStart      uint64          `db:"f_block_start"                  json:"blockStart"`
	BlockCurrent    uint64          `db:"f_block_current"                json:"blockCurrent,omitempty"`
	BlockEnd        uint64          `db:"f_block_end,default='0'"        json:"blockEnd,omitempty"`
	Topic0          string          `db:"f_topic0,default=''"            json:"topic0,omitempty"`
	Topic1          string          `db:"f_topic1,default=''"            json:"topic1,omitempty"`
	Topic2          string          `db:"f_topic2,default=''"            json:"topic2,omitempty"`
	Topic3          string          `db:"f_topic3,default=''"            json:"topic3,omitempty"`
}
