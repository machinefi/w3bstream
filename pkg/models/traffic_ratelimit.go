package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
)

// TrafficRateLimit traffic rate limit for each project
// @def primary                       ID
// @def unique_index UI_ratelimit_id  RateLimitID
// @def unique_index UI_prj_api_type  ProjectID ApiType
//
//go:generate toolkit gen model TrafficRateLimit --database DB
type TrafficRateLimit struct {
	datatypes.PrimaryID
	RelRateLimit
	RelProject
	RateLimitInfo
	datatypes.OperationTimesWithDeleted
}

type RelRateLimit struct {
	RateLimitID types.SFID `db:"f_ratelimit_id" json:"rateLimitID"`
}

type RateLimitInfo struct {
	Threshold int                    `db:"f_threshold"                    json:"threshold"`
	CycleNum  int                    `db:"f_cycle_num"                    json:"cycleNum"`
	CycleUnit enums.TrafficCycle     `db:"f_cycle_unit"                    json:"cycleUnit"`
	ApiType   enums.RateLimitApiType `db:"f_api_type"                      json:"apiType"`
}
