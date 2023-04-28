package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

// CronJob schema for cron job information
// @def primary                      ID
// @def unique_index UI_cron_job_id  CronJobID
// @def unique_index UI_cron         ProjectID CronExpressions EventType
//
//go:generate toolkit gen model CronJob --database DB
type CronJob struct {
	datatypes.PrimaryID
	RelCronJob
	RelProject
	CronJobInfo
	datatypes.OperationTimesWithDeleted
}

type RelCronJob struct {
	CronJobID types.SFID `db:"f_cron_job_id" json:"cronJobID"`
}

type CronJobInfo struct {
	CronExpressions string `db:"f_cron_expressions"    json:"cronExpressions"`
	EventType       string `db:"f_event_type"          json:"eventType,omitempty"`
}
