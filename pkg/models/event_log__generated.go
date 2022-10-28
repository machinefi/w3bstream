// This is a generated source file. DO NOT EDIT
// Source: models/event_log__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/sqlx"
	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"
)

var EventLogTable *builder.Table

func init() {
	EventLogTable = DB.Register(&EventLog{})
}

type EventLogIterator struct {
}

func (*EventLogIterator) New() interface{} {
	return &EventLog{}
}

func (*EventLogIterator) Resolve(v interface{}) *EventLog {
	return v.(*EventLog)
}

func (*EventLog) TableName() string {
	return "t_event_log"
}

func (*EventLog) TableDesc() []string {
	return []string{
		"EventLog database model event",
	}
}

func (*EventLog) Comments() map[string]string {
	return map[string]string{}
}

func (*EventLog) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*EventLog) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*EventLog) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (*EventLog) Indexes() builder.Indexes {
	return builder.Indexes{
		"i_applet_id": []string{
			"ProjectID",
		},
		"i_project_id": []string{
			"ProjectID",
		},
		"i_publisher_id": []string{
			"PublisherID",
		},
	}
}

func (m *EventLog) IndexFieldNames() []string {
	return []string{
		"EventID",
		"ID",
		"ProjectID",
		"PublisherID",
	}
}

func (*EventLog) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_event_id": []string{
			"EventID",
		},
	}
}

func (*EventLog) UniqueIndexUIEventID() string {
	return "ui_event_id"
}

func (m *EventLog) ColID() *builder.Column {
	return EventLogTable.ColByFieldName(m.FieldID())
}

func (*EventLog) FieldID() string {
	return "ID"
}

func (m *EventLog) ColEventID() *builder.Column {
	return EventLogTable.ColByFieldName(m.FieldEventID())
}

func (*EventLog) FieldEventID() string {
	return "EventID"
}

func (m *EventLog) ColProjectID() *builder.Column {
	return EventLogTable.ColByFieldName(m.FieldProjectID())
}

func (*EventLog) FieldProjectID() string {
	return "ProjectID"
}

func (m *EventLog) ColAppletID() *builder.Column {
	return EventLogTable.ColByFieldName(m.FieldAppletID())
}

func (*EventLog) FieldAppletID() string {
	return "AppletID"
}

func (m *EventLog) ColPublisherID() *builder.Column {
	return EventLogTable.ColByFieldName(m.FieldPublisherID())
}

func (*EventLog) FieldPublisherID() string {
	return "PublisherID"
}

func (m *EventLog) ColCreatedAt() *builder.Column {
	return EventLogTable.ColByFieldName(m.FieldCreatedAt())
}

func (*EventLog) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *EventLog) ColUpdatedAt() *builder.Column {
	return EventLogTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*EventLog) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *EventLog) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *EventLog) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *EventLog) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]EventLog, error) {
	var (
		tbl = db.T(m)
		lst = make([]EventLog, 0)
	)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("EventLog.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *EventLog) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("EventLog.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *EventLog) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("EventLog.FetchByID"),
			),
		m,
	)
	return err
}

func (m *EventLog) FetchByEventID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("EventID").Eq(m.EventID),
					),
				),
				builder.Comment("EventLog.FetchByEventID"),
			),
		m,
	)
	return err
}

func (m *EventLog) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ID").Eq(m.ID),
				),
				builder.Comment("EventLog.UpdateByIDWithFVs"),
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

func (m *EventLog) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *EventLog) UpdateByEventIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("EventID").Eq(m.EventID),
				),
				builder.Comment("EventLog.UpdateByEventIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByEventID(db)
	}
	return nil
}

func (m *EventLog) UpdateByEventID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByEventIDWithFVs(db, fvs)
}

func (m *EventLog) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("EventLog.Delete"),
			),
	)
	return err
}

func (m *EventLog) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("EventLog.DeleteByID"),
			),
	)
	return err
}

func (m *EventLog) DeleteByEventID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("EventID").Eq(m.EventID),
					),
				),
				builder.Comment("EventLog.DeleteByEventID"),
			),
	)
	return err
}
