package models

import (
	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/sqlx/datatypes"
)

// EventLog database model event
// @def primary                     ID
// @def unique_index UI_event_id    EventID
// @def index        I_project_id   ProjectID
// @def index        I_applet_id    ProjectID
// @def index        I_publisher_id PublisherID
//
//go:generate toolkit gen model EventLog --database DB
type EventLog struct {
	datatypes.PrimaryID
	RefEventID
	EventInfo
	datatypes.OperationTimes
}

type RefEventID struct {
	EventID types.SFID `db:"f_event_id" json:"eventID"`
}

type EventInfo struct {
	RelProject
	RelApplet
	RelPublisher
}
