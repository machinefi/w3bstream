// This is a generated source file. DO NOT EDIT
// Source: models/strategy__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var StrategyTable *builder.Table

func init() {
	StrategyTable = DB.Register(&Strategy{})
}

type StrategyIterator struct {
}

func (*StrategyIterator) New() interface{} {
	return &Strategy{}
}

func (*StrategyIterator) Resolve(v interface{}) *Strategy {
	return v.(*Strategy)
}

func (*Strategy) TableName() string {
	return "t_strategy"
}

func (*Strategy) TableDesc() []string {
	return []string{
		"Strategy event route strategy",
	}
}

func (*Strategy) Comments() map[string]string {
	return map[string]string{
		"AutoCollectMetric": "AutoCollectMetric if allow host collect event data for metering",
		"EventType":         "EventType user defined event type",
		"Handler":           "Handler wasm handler fn name",
	}
}

func (*Strategy) ColDesc() map[string][]string {
	return map[string][]string{
		"AutoCollectMetric": []string{
			"AutoCollectMetric if allow host collect event data for metering",
		},
		"EventType": []string{
			"EventType user defined event type",
		},
		"Handler": []string{
			"Handler wasm handler fn name",
		},
	}
}

func (*Strategy) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*Strategy) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Strategy) IndexFieldNames() []string {
	return []string{
		"AppletID",
		"EventType",
		"ID",
		"ProjectID",
		"StrategyID",
	}
}

func (*Strategy) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_prj_app_event": []string{
			"ProjectID",
			"AppletID",
			"EventType",
			"DeletedAt",
		},
		"ui_strategy_id": []string{
			"StrategyID",
			"DeletedAt",
		},
	}
}

func (*Strategy) UniqueIndexUIPrjAppEvent() string {
	return "ui_prj_app_event"
}

func (*Strategy) UniqueIndexUIStrategyID() string {
	return "ui_strategy_id"
}

func (m *Strategy) ColID() *builder.Column {
	return StrategyTable.ColByFieldName(m.FieldID())
}

func (*Strategy) FieldID() string {
	return "ID"
}

func (m *Strategy) ColStrategyID() *builder.Column {
	return StrategyTable.ColByFieldName(m.FieldStrategyID())
}

func (*Strategy) FieldStrategyID() string {
	return "StrategyID"
}

func (m *Strategy) ColProjectID() *builder.Column {
	return StrategyTable.ColByFieldName(m.FieldProjectID())
}

func (*Strategy) FieldProjectID() string {
	return "ProjectID"
}

func (m *Strategy) ColAppletID() *builder.Column {
	return StrategyTable.ColByFieldName(m.FieldAppletID())
}

func (*Strategy) FieldAppletID() string {
	return "AppletID"
}

func (m *Strategy) ColEventType() *builder.Column {
	return StrategyTable.ColByFieldName(m.FieldEventType())
}

func (*Strategy) FieldEventType() string {
	return "EventType"
}

func (m *Strategy) ColHandler() *builder.Column {
	return StrategyTable.ColByFieldName(m.FieldHandler())
}

func (*Strategy) FieldHandler() string {
	return "Handler"
}

func (m *Strategy) ColAutoCollectMetric() *builder.Column {
	return StrategyTable.ColByFieldName(m.FieldAutoCollectMetric())
}

func (*Strategy) FieldAutoCollectMetric() string {
	return "AutoCollectMetric"
}

func (m *Strategy) ColCreatedAt() *builder.Column {
	return StrategyTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Strategy) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Strategy) ColUpdatedAt() *builder.Column {
	return StrategyTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Strategy) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Strategy) ColDeletedAt() *builder.Column {
	return StrategyTable.ColByFieldName(m.FieldDeletedAt())
}

func (*Strategy) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *Strategy) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
	var (
		tbl  = db.T(m)
		fvs  = builder.FieldValueFromStructByNoneZero(m)
		cond = []builder.SqlCondition{tbl.ColByFieldName("DeletedAt").Eq(0)}
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

func (m *Strategy) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Strategy) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Strategy, error) {
	var (
		tbl = db.T(m)
		lst = make([]Strategy, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Strategy.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Strategy) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Strategy.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Strategy) FetchByID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ID").Eq(m.ID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Strategy.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Strategy) FetchByProjectIDAndAppletIDAndEventType(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
						tbl.ColByFieldName("AppletID").Eq(m.AppletID),
						tbl.ColByFieldName("EventType").Eq(m.EventType),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Strategy.FetchByProjectIDAndAppletIDAndEventType"),
			),
		m,
	)
	return err
}

func (m *Strategy) FetchByStrategyID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("StrategyID").Eq(m.StrategyID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Strategy.FetchByStrategyID"),
			),
		m,
	)
	return err
}

func (m *Strategy) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ID").Eq(m.ID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Strategy.UpdateByIDWithFVs"),
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

func (m *Strategy) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Strategy) UpdateByProjectIDAndAppletIDAndEventTypeWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
					tbl.ColByFieldName("AppletID").Eq(m.AppletID),
					tbl.ColByFieldName("EventType").Eq(m.EventType),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Strategy.UpdateByProjectIDAndAppletIDAndEventTypeWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByProjectIDAndAppletIDAndEventType(db)
	}
	return nil
}

func (m *Strategy) UpdateByProjectIDAndAppletIDAndEventType(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByProjectIDAndAppletIDAndEventTypeWithFVs(db, fvs)
}

func (m *Strategy) UpdateByStrategyIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("StrategyID").Eq(m.StrategyID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Strategy.UpdateByStrategyIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByStrategyID(db)
	}
	return nil
}

func (m *Strategy) UpdateByStrategyID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByStrategyIDWithFVs(db, fvs)
}

func (m *Strategy) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Strategy.Delete"),
			),
	)
	return err
}

func (m *Strategy) DeleteByID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ID").Eq(m.ID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Strategy.DeleteByID"),
			),
	)
	return err
}

func (m *Strategy) SoftDeleteByID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	fvs := builder.FieldValues{}

	if _, ok := fvs["DeletedAt"]; !ok {
		fvs["DeletedAt"] = types.Timestamp{Time: time.Now()}
	}

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	_, err := db.Exec(
		builder.Update(db.T(m)).
			Where(
				builder.And(
					tbl.ColByFieldName("ID").Eq(m.ID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Strategy.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *Strategy) DeleteByProjectIDAndAppletIDAndEventType(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
						tbl.ColByFieldName("AppletID").Eq(m.AppletID),
						tbl.ColByFieldName("EventType").Eq(m.EventType),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Strategy.DeleteByProjectIDAndAppletIDAndEventType"),
			),
	)
	return err
}

func (m *Strategy) SoftDeleteByProjectIDAndAppletIDAndEventType(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	fvs := builder.FieldValues{}

	if _, ok := fvs["DeletedAt"]; !ok {
		fvs["DeletedAt"] = types.Timestamp{Time: time.Now()}
	}

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	_, err := db.Exec(
		builder.Update(db.T(m)).
			Where(
				builder.And(
					tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
					tbl.ColByFieldName("AppletID").Eq(m.AppletID),
					tbl.ColByFieldName("EventType").Eq(m.EventType),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Strategy.SoftDeleteByProjectIDAndAppletIDAndEventType"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *Strategy) DeleteByStrategyID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("StrategyID").Eq(m.StrategyID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Strategy.DeleteByStrategyID"),
			),
	)
	return err
}

func (m *Strategy) SoftDeleteByStrategyID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	fvs := builder.FieldValues{}

	if _, ok := fvs["DeletedAt"]; !ok {
		fvs["DeletedAt"] = types.Timestamp{Time: time.Now()}
	}

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	_, err := db.Exec(
		builder.Update(db.T(m)).
			Where(
				builder.And(
					tbl.ColByFieldName("StrategyID").Eq(m.StrategyID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Strategy.SoftDeleteByStrategyID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
