// This is a generated source file. DO NOT EDIT
// Source: models/event__generated.go

package models

import (
	"fmt"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var EventTable *builder.Table

func init() {
	EventTable = DB.Register(&Event{})
}

type EventIterator struct {
}

func (*EventIterator) New() interface{} {
	return &Event{}
}

func (*EventIterator) Resolve(v interface{}) *Event {
	return v.(*Event)
}

func (*Event) TableName() string {
	return "t_event"
}

func (*Event) TableDesc() []string {
	return []string{
		"EventLog database model for event context",
	}
}

func (*Event) Comments() map[string]string {
	return map[string]string{
		"AccountID":    "AccountID account ID",
		"CompletedAt":  "CompletedAt event completed timestamp(epoch milliseconds)",
		"Error":        "Error wasm handle error message",
		"EventID":      "EventID event ID",
		"EventType":    "EventType event type",
		"From":         "From channel type: MQTT HTTP",
		"HandledAt":    "HandledAt event handled timestamp(epoch milliseconds)",
		"Handler":      "Handler wasm exported handling func",
		"Index":        "Index strategy index for different handler",
		"Input":        "Input event payload",
		"InstanceID":   "InstanceID instance ID",
		"ProjectID":    "ProjectID project ID",
		"ProjectName":  "ProjectName project name",
		"PublishedAt":  "PublishedAt event published timestamp(epoch milliseconds) from request",
		"PublisherID":  "PublisherID publisher ID",
		"PublisherKey": "PublisherKey publisher key",
		"ReceivedAt":   "ReceivedAt event received timestamp(epoch milliseconds)",
		"ResultCode":   "ResultCode wasm handle result code",
		"Stage":        "Stage event handle stage: RECEIVED, HANDLED and COMPLETED",
		"Total":        "Total strategy total",
	}
}

func (*Event) ColDesc() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"AccountID account ID",
		},
		"CompletedAt": []string{
			"CompletedAt event completed timestamp(epoch milliseconds)",
		},
		"Error": []string{
			"Error wasm handle error message",
		},
		"EventID": []string{
			"EventID event ID",
		},
		"EventType": []string{
			"EventType event type",
		},
		"From": []string{
			"From channel type: MQTT HTTP",
		},
		"HandledAt": []string{
			"HandledAt event handled timestamp(epoch milliseconds)",
		},
		"Handler": []string{
			"Handler wasm exported handling func",
		},
		"Index": []string{
			"Index strategy index for different handler",
		},
		"Input": []string{
			"Input event payload",
		},
		"InstanceID": []string{
			"InstanceID instance ID",
		},
		"ProjectID": []string{
			"ProjectID project ID",
		},
		"ProjectName": []string{
			"ProjectName project name",
		},
		"PublishedAt": []string{
			"PublishedAt event published timestamp(epoch milliseconds) from request",
		},
		"PublisherID": []string{
			"PublisherID publisher ID",
		},
		"PublisherKey": []string{
			"PublisherKey publisher key",
		},
		"ReceivedAt": []string{
			"ReceivedAt event received timestamp(epoch milliseconds)",
		},
		"ResultCode": []string{
			"ResultCode wasm handle result code",
		},
		"Stage": []string{
			"Stage event handle stage: RECEIVED, HANDLED and COMPLETED",
		},
		"Total": []string{
			"Total strategy total",
		},
	}
}

func (*Event) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*Event) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (*Event) Indexes() builder.Indexes {
	return builder.Indexes{
		"i_completed_at": []string{
			"CompletedAt",
		},
		"i_event_id": []string{
			"EventID",
		},
		"i_event_type": []string{
			"EventType",
		},
		"i_handled_at": []string{
			"HandledAt",
		},
		"i_handler": []string{
			"Handler",
		},
		"i_instance_id": []string{
			"InstanceID",
		},
		"i_project_id": []string{
			"ProjectID",
		},
		"i_published_at": []string{
			"PublishedAt",
		},
		"i_publisher_id": []string{
			"PublisherID",
		},
		"i_received_at": []string{
			"ReceivedAt",
		},
		"i_result_code": []string{
			"ResultCode",
		},
	}
}

func (m *Event) IndexFieldNames() []string {
	return []string{
		"CompletedAt",
		"EventID",
		"EventType",
		"HandledAt",
		"Handler",
		"ID",
		"InstanceID",
		"ProjectID",
		"PublishedAt",
		"PublisherID",
		"ReceivedAt",
		"ResultCode",
	}
}

func (m *Event) ColID() *builder.Column {
	return EventTable.ColByFieldName(m.FieldID())
}

func (*Event) FieldID() string {
	return "ID"
}

func (m *Event) ColStage() *builder.Column {
	return EventTable.ColByFieldName(m.FieldStage())
}

func (*Event) FieldStage() string {
	return "Stage"
}

func (m *Event) ColFrom() *builder.Column {
	return EventTable.ColByFieldName(m.FieldFrom())
}

func (*Event) FieldFrom() string {
	return "From"
}

func (m *Event) ColAccountID() *builder.Column {
	return EventTable.ColByFieldName(m.FieldAccountID())
}

func (*Event) FieldAccountID() string {
	return "AccountID"
}

func (m *Event) ColProjectID() *builder.Column {
	return EventTable.ColByFieldName(m.FieldProjectID())
}

func (*Event) FieldProjectID() string {
	return "ProjectID"
}

func (m *Event) ColProjectName() *builder.Column {
	return EventTable.ColByFieldName(m.FieldProjectName())
}

func (*Event) FieldProjectName() string {
	return "ProjectName"
}

func (m *Event) ColPublisherID() *builder.Column {
	return EventTable.ColByFieldName(m.FieldPublisherID())
}

func (*Event) FieldPublisherID() string {
	return "PublisherID"
}

func (m *Event) ColPublisherKey() *builder.Column {
	return EventTable.ColByFieldName(m.FieldPublisherKey())
}

func (*Event) FieldPublisherKey() string {
	return "PublisherKey"
}

func (m *Event) ColEventID() *builder.Column {
	return EventTable.ColByFieldName(m.FieldEventID())
}

func (*Event) FieldEventID() string {
	return "EventID"
}

func (m *Event) ColIndex() *builder.Column {
	return EventTable.ColByFieldName(m.FieldIndex())
}

func (*Event) FieldIndex() string {
	return "Index"
}

func (m *Event) ColTotal() *builder.Column {
	return EventTable.ColByFieldName(m.FieldTotal())
}

func (*Event) FieldTotal() string {
	return "Total"
}

func (m *Event) ColEventType() *builder.Column {
	return EventTable.ColByFieldName(m.FieldEventType())
}

func (*Event) FieldEventType() string {
	return "EventType"
}

func (m *Event) ColInstanceID() *builder.Column {
	return EventTable.ColByFieldName(m.FieldInstanceID())
}

func (*Event) FieldInstanceID() string {
	return "InstanceID"
}

func (m *Event) ColHandler() *builder.Column {
	return EventTable.ColByFieldName(m.FieldHandler())
}

func (*Event) FieldHandler() string {
	return "Handler"
}

func (m *Event) ColInput() *builder.Column {
	return EventTable.ColByFieldName(m.FieldInput())
}

func (*Event) FieldInput() string {
	return "Input"
}

func (m *Event) ColResultCode() *builder.Column {
	return EventTable.ColByFieldName(m.FieldResultCode())
}

func (*Event) FieldResultCode() string {
	return "ResultCode"
}

func (m *Event) ColError() *builder.Column {
	return EventTable.ColByFieldName(m.FieldError())
}

func (*Event) FieldError() string {
	return "Error"
}

func (m *Event) ColPublishedAt() *builder.Column {
	return EventTable.ColByFieldName(m.FieldPublishedAt())
}

func (*Event) FieldPublishedAt() string {
	return "PublishedAt"
}

func (m *Event) ColReceivedAt() *builder.Column {
	return EventTable.ColByFieldName(m.FieldReceivedAt())
}

func (*Event) FieldReceivedAt() string {
	return "ReceivedAt"
}

func (m *Event) ColHandledAt() *builder.Column {
	return EventTable.ColByFieldName(m.FieldHandledAt())
}

func (*Event) FieldHandledAt() string {
	return "HandledAt"
}

func (m *Event) ColCompletedAt() *builder.Column {
	return EventTable.ColByFieldName(m.FieldCompletedAt())
}

func (*Event) FieldCompletedAt() string {
	return "CompletedAt"
}

func (m *Event) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
	var (
		tbl  = db.T(m)
		fvs  = builder.FieldValueFromStructByNoneZero(m)
		cond = make([]builder.SqlCondition, 0)
	)

	for _, fn := range m.IndexFieldNames() {
		if v, ok := fvs[fn]; ok {
			cond = append(cond, tbl.ColByFieldName(fn).Eq(v))
			delete(fvs, fn)
		}
	}
	if len(cond) == 0 {
		panic(fmt.Errorf("no field for indexes has value"))
	}
	for fn, v := range fvs {
		cond = append(cond, tbl.ColByFieldName(fn).Eq(v))
	}
	return builder.And(cond...)
}

func (m *Event) Create(db sqlx.DBExecutor) error {

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Event) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Event, error) {
	var (
		tbl = db.T(m)
		lst = make([]Event, 0)
	)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Event.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Event) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Event.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Event) FetchByID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ID").Eq(m.ID),
					),
				),
				builder.Comment("Event.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Event) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ID").Eq(m.ID),
				),
				builder.Comment("Event.UpdateByIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByID(db)
	}
	return nil
}

func (m *Event) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Event) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Event.Delete"),
			),
	)
	return err
}

func (m *Event) DeleteByID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ID").Eq(m.ID),
					),
				),
				builder.Comment("Event.DeleteByID"),
			),
	)
	return err
}
