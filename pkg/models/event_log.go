package models

import "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"

// EventLog database model event
// @def primary              ID
// @def index I_event_id     EventID
// @def index I_project_id   ProjectID
// @def index I_applet_id    ProjectID
// @def index I_publisher_id PublisherID
//
//go:generate toolkit gen model EventLog --database DB
type EventLog struct {
	datatypes.PrimaryID
	EventInfo
	datatypes.OperationTimes
}

type EventInfo struct {
	EventID string `db:"f_event_id" json:"eventID"`
	RelProject
	RelPublisher
	// PublishedAt the timestamp when device publish event
	PublishedAt int64 `db:"f_published_at" json:"publishedAt"`
	// ReceivedAt the timestamp when event received by us
	ReceivedAt int64 `db:"f_received_at" json:"receivedAt"`
	// RespondedAt the timestamp when event handled and send response
	RespondedAt int64 `db:"f_responded_at" json:"respondedAt"`
}
