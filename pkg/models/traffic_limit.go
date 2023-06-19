package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
)

// TrafficLimit traffic limit for each project
// @def primary                           ID
// @def unique_index UI_traffic_limit_id  TrafficLimitID
// @def unique_index UI_prj_api_type      ProjectID ApiType
//
//go:generate toolkit gen model TrafficLimit --database DB
type TrafficLimit struct {
	datatypes.PrimaryID
	RelTrafficLimit
	RelProject
	TrafficLimitInfo
	datatypes.OperationTimesWithDeleted
}

type RelTrafficLimit struct {
	TrafficLimitID types.SFID `db:"f_traffic_limit_id" json:"trafficLimitID"`
}

type TrafficLimitInfo struct {
	Threshold int                    `db:"f_threshold"                     json:"threshold"`
	Duration  types.Duration         `db:"f_duration"                      json:"duration"`
	ApiType   enums.TrafficLimitType `db:"f_api_type"                      json:"apiType"`
	StartAt   types.Timestamp        `db:"f_start_at,default='0'"          json:"startAt"`
}
