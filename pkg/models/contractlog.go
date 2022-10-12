package models

import (
	"github.com/iotexproject/Bumblebee/kit/sqlx/datatypes"
)

// Contractlog database model contractlog
// @def primary                   ID
//
//go:generate toolkit gen model Contractlog --database DB
type Contractlog struct {
	datatypes.PrimaryID
	RelContractlog
	ContractlogInfo
	datatypes.OperationTimes
}

type RelContractlog struct {
	ContractlogID string `db:"f_contractlog_id" json:"contractlogID"`
}

type ContractlogInfo struct {
	ProjectName     string `db:"f_project_name"                 json:"projectName"`
	EventType       string `db:"f_event_type"                   json:"eventType"`
	ChainID         uint64 `db:"f_chainID"                      json:"chainID"`
	ContractAddress string `db:"f_contractAddress"              json:"contractAddress"`
	BlockStart      uint64 `db:"f_blockStart"                   json:"blockStart"`
	BlockCurrent    uint64 `db:"f_blockCurrent"                 json:"blockCurrent"`
	BlockEnd        uint64 `db:"f_blockEnd,default='0'"         json:"blockEnd,omitempty"`
	Topic1          string `db:"f_topic1,default=''"            json:"topic1,omitempty"`
	Topic2          string `db:"f_topic2,default=''"            json:"topic2,omitempty"`
	Topic3          string `db:"f_topic3,default=''"            json:"topic3,omitempty"`
}
