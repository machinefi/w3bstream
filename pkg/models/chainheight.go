package models

import (
	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/sqlx/datatypes"

	"github.com/iotexproject/w3bstream/pkg/enums"
)

// ChainHeight database model chainheight
// @def primary                   ID
//
//go:generate toolkit gen model ChainHeight --database DB
type ChainHeight struct {
	datatypes.PrimaryID
	RelChainHeight
	ChainHeightData
	datatypes.OperationTimes
}

type RelChainHeight struct {
	ChainHeightID types.SFID `db:"f_chain_height_id" json:"chainHeightID"`
}

type ChainHeightData struct {
	ProjectName string `db:"f_project_name"                 json:"projectName"`
	Finished    bool   `db:"f_finished,default='false'"     json:"finished,omitempty"`
	ChainHeightInfo
}

type ChainHeightInfo struct {
	EventType enums.EventType `db:"f_event_type"                   json:"eventType,omitempty"`
	ChainID   uint64          `db:"f_chain_id"                     json:"chainID"`
	Height    uint64          `db:"f_height"                       json:"height"`
}
