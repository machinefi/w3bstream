package models

import (
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
)

// EventLog database model for event context
// @def primary              ID
// @def index I_event_id     EventID
// @def index I_event_type   EventType
// @def index I_project_id   ProjectID
// @def index I_publisher_id PublisherID
// @def index I_instance_id  InstanceID
// @def index I_handler      Handler
// @def index I_result_code  ResultCode
// @def index I_published_at PublishedAt
// @def index I_received_at  ReceivedAt
// @def index I_handled_at   HandledAt
// @def index I_completed_at CompletedAt
//
//go:generate toolkit gen model Event --database DB
type Event struct {
	datatypes.PrimaryID
	EventContext
}

type EventContext struct {
	// Stage event handle stage: RECEIVED, HANDLED and COMPLETED
	Stage enums.EventStage `db:"f_stage" json:"stage"`
	// From channel type: MQTT HTTP
	From enums.EventSource `db:"f_from" json:"from"`
	// AccountID account ID
	AccountID types.SFID `db:"f_account_id" json:"accountID"`
	// ProjectID project ID
	ProjectID types.SFID `db:"f_project_id" json:"projectID"`
	// ProjectName project name
	ProjectName string `db:"f_project_name" json:"projectName"`
	// PublisherID publisher ID
	PublisherID types.SFID `db:"f_publisher_id" json:"publisherID"`
	// PublisherKey publisher key
	PublisherKey string `db:"f_publisher_key" json:"publisherKey"`
	// EventID event ID
	EventID string `db:"f_event_id" json:"eventID"`
	// Index strategy index for different handler
	Index int `db:"f_index,default='0'" json:"index"`
	// Total strategy total
	Total int `db:"f_total,default='1'" json:"total"`
	// EventType event type
	EventType string `db:"f_event_type" json:"eventType"`
	// InstanceID instance ID
	InstanceID types.SFID `db:"f_instance_id" json:"instanceID"`
	// Handler wasm exported handling func
	Handler string `db:"f_handler" json:"handler"`
	// Input event payload
	Input []byte `db:"f_input" json:"input"`
	// ResultCode wasm handle result code
	ResultCode int32 `db:"f_result_code,default='0'" json:"resultCode"`
	// Error wasm handle error message
	Error string `db:"f_error,default=''"   json:"error"`
	// PublishedAt event published timestamp(epoch milliseconds) from request
	PublishedAt int64 `db:"f_published_at" json:"publishedAt"`
	// ReceivedAt event received timestamp(epoch milliseconds)
	ReceivedAt int64 `db:"f_received_at"  json:"receivedAt"`
	// HandledAt event handled timestamp(epoch milliseconds)
	HandledAt int64 `db:"f_handled_at,default='0'"   json:"handledAt"`
	// CompletedAt event completed timestamp(epoch milliseconds)
	CompletedAt int64 `db:"f_completed_at,default='0'" json:"completedAt"`
	// AutoCollect if do geo collection
	AutoCollect datatypes.Bool `db:"f_auto_collection,default='2'" json:"autoCollection"`
}

func BatchFetchEvents(d sqlx.DBExecutor, adds ...builder.Addition) (results []*Event, err error) {
	m := &Event{}
	t := d.T(m)

	err = d.QueryAndScan(builder.Select(nil).From(t, adds...), &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func BatchFetchLastUnhandledEvents(d sqlx.DBExecutor, batch int64, prj types.SFID) ([]*Event, error) {
	m := &Event{}
	return BatchFetchEvents(
		d,
		builder.Where(
			builder.And(
				m.ColStage().Eq(enums.EVENT_STAGE__RECEIVED),
				m.ColProjectID().Eq(prj),
			),
		),
		builder.OrderBy(
			builder.AscOrder(m.ColReceivedAt()),
		),
		builder.Limit(batch),
		builder.Comment("BatchFetchLastUnhandledEvents"),
	)
}

func BatchCreateEvents(d sqlx.DBExecutor, vs ...*Event) error {
	m := &Event{}

	cols := reflect.TypeOf(m).Elem().NumField()
	args := make([]any, 0, len(vs)*cols)
	for _, v := range vs {
		args = append(args,
			v.Stage,
			v.From,
			v.AccountID,
			v.ProjectID,
			v.ProjectName,
			v.PublisherID,
			v.PublisherKey,
			v.EventID,
			v.Index,
			v.Total,
			v.EventType,
			v.InstanceID,
			v.Handler,
			v.Input,
			v.ResultCode,
			v.Error,
			v.PublishedAt,
			v.ReceivedAt,
			v.HandledAt,
			v.CompletedAt,
			v.AutoCollect,
		)
	}
	if len(args) == 0 {
		return nil
	}

	t := d.T(m)
	_, err := d.Exec(builder.Insert().Into(t).Values(
		builder.Cols(
			m.ColStage().Name,
			m.ColFrom().Name,
			m.ColAccountID().Name,
			m.ColProjectID().Name,
			m.ColProjectName().Name,
			m.ColPublisherID().Name,
			m.ColPublisherKey().Name,
			m.ColEventID().Name,
			m.ColIndex().Name,
			m.ColTotal().Name,
			m.ColEventType().Name,
			m.ColInstanceID().Name,
			m.ColHandler().Name,
			m.ColInput().Name,
			m.ColResultCode().Name,
			m.ColError().Name,
			m.ColPublishedAt().Name,
			m.ColReceivedAt().Name,
			m.ColHandledAt().Name,
			m.ColCompletedAt().Name,
			m.ColAutoCollect().Name,
		),
		args...,
	))
	return err
}
