package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
)

// Strategy event route strategy
// @def primary                       ID
// @def unique_index UI_strategy_id   StrategyID
// @def unique_index UI_prj_app_event ProjectID AppletID EventType
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
	EventType string `db:"f_event_type" json:"eventType"`
	Handler   string `db:"f_handler"    json:"handler"`
}

var DefaultStrategyInfo = StrategyInfo{
	EventType: enums.EVENTTYPEDEFAULT,
	Handler:   "start",
}
