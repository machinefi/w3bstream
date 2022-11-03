package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

// ChainHeight database model chainheight
// @def primary                   ID
// @def unique_index UI_chain_height_id   ChainHeightID
// @def unique_index UI_chain_height_uniq ProjectName EventType ChainID Height Uniq
//
//go:generate toolkit gen model ChainHeight --database MonitorDB
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
	ProjectName string         `db:"f_project_name"                 json:"projectName"`
	Finished    datatypes.Bool `db:"f_finished,default='2'"         json:"-"`
	Uniq        types.SFID     `db:"f_uniq,default='0'"             json:"-"`
	ChainHeightInfo
}

type ChainHeightInfo struct {
	EventType string `db:"f_event_type"                   json:"eventType,omitempty"`
	ChainID   uint64 `db:"f_chain_id"                     json:"chainID"`
	Height    uint64 `db:"f_height"                       json:"height"`
}
