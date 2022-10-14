package models

import (
	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/sqlx/datatypes"

	"github.com/iotexproject/w3bstream/pkg/enums"
)

// Strategy event route strategy
// @def primary                       ID
// @def unique_index UI_strategy_id   StrategyID
// @def unique_index UI_prj_app_event ProjectID AppletID EventType Handler
//
//go:generate toolkit gen model Strategy --database DB
type Strategy struct {
	datatypes.PrimaryID
	RelStrategy
	RelProject
	RelApplet
	StrategyInfo
	datatypes.OperationTimesWithDeleted
}

type RelStrategy struct {
	StrategyID types.SFID `db:"f_strategy_id" json:"strategyID"`
}

type StrategyInfo struct {
	EventType enums.EventType `db:"f_event_type" json:"eventType"`
	Handler   string          `db:"f_handler"    json:"handler"`
}

var DefaultStrategyInfo = StrategyInfo{
	EventType: enums.EVENT_TYPE__ANY,
	Handler:   "start",
}
