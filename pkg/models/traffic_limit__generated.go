// This is a generated source file. DO NOT EDIT
// Source: models/traffic_limit__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var TrafficLimitTable *builder.Table

func init() {
	TrafficLimitTable = DB.Register(&TrafficLimit{})
}

type TrafficLimitIterator struct {
}

func (*TrafficLimitIterator) New() interface{} {
	return &TrafficLimit{}
}

func (*TrafficLimitIterator) Resolve(v interface{}) *TrafficLimit {
	return v.(*TrafficLimit)
}

func (*TrafficLimit) TableName() string {
	return "t_traffic_limit"
}

func (*TrafficLimit) TableDesc() []string {
	return []string{
		"TrafficLimit traffic limit for each project",
	}
}

func (*TrafficLimit) Comments() map[string]string {
	return map[string]string{}
}

func (*TrafficLimit) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*TrafficLimit) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*TrafficLimit) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *TrafficLimit) IndexFieldNames() []string {
	return []string{
		"ApiType",
		"ID",
		"ProjectID",
		"TrafficLimitID",
	}
}

func (*TrafficLimit) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_prj_api_type": []string{
			"ProjectID",
			"ApiType",
			"DeletedAt",
		},
		"ui_traffic_limit_id": []string{
			"TrafficLimitID",
			"DeletedAt",
		},
	}
}

func (*TrafficLimit) UniqueIndexUIPrjAPIType() string {
	return "ui_prj_api_type"
}

func (*TrafficLimit) UniqueIndexUITrafficLimitID() string {
	return "ui_traffic_limit_id"
}

func (m *TrafficLimit) ColID() *builder.Column {
	return TrafficLimitTable.ColByFieldName(m.FieldID())
}

func (*TrafficLimit) FieldID() string {
	return "ID"
}

func (m *TrafficLimit) ColTrafficLimitID() *builder.Column {
	return TrafficLimitTable.ColByFieldName(m.FieldTrafficLimitID())
}

func (*TrafficLimit) FieldTrafficLimitID() string {
	return "TrafficLimitID"
}

func (m *TrafficLimit) ColProjectID() *builder.Column {
	return TrafficLimitTable.ColByFieldName(m.FieldProjectID())
}

func (*TrafficLimit) FieldProjectID() string {
	return "ProjectID"
}

func (m *TrafficLimit) ColThreshold() *builder.Column {
	return TrafficLimitTable.ColByFieldName(m.FieldThreshold())
}

func (*TrafficLimit) FieldThreshold() string {
	return "Threshold"
}

func (m *TrafficLimit) ColDuration() *builder.Column {
	return TrafficLimitTable.ColByFieldName(m.FieldDuration())
}

func (*TrafficLimit) FieldDuration() string {
	return "Duration"
}

func (m *TrafficLimit) ColApiType() *builder.Column {
	return TrafficLimitTable.ColByFieldName(m.FieldApiType())
}

func (*TrafficLimit) FieldApiType() string {
	return "ApiType"
}

func (m *TrafficLimit) ColCreatedAt() *builder.Column {
	return TrafficLimitTable.ColByFieldName(m.FieldCreatedAt())
}

func (*TrafficLimit) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *TrafficLimit) ColUpdatedAt() *builder.Column {
	return TrafficLimitTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*TrafficLimit) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *TrafficLimit) ColDeletedAt() *builder.Column {
	return TrafficLimitTable.ColByFieldName(m.FieldDeletedAt())
}

func (*TrafficLimit) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *TrafficLimit) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *TrafficLimit) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *TrafficLimit) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]TrafficLimit, error) {
	var (
		tbl = db.T(m)
		lst = make([]TrafficLimit, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("TrafficLimit.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *TrafficLimit) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("TrafficLimit.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *TrafficLimit) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("TrafficLimit.FetchByID"),
			),
		m,
	)
	return err
}

func (m *TrafficLimit) FetchByProjectIDAndApiType(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
						tbl.ColByFieldName("ApiType").Eq(m.ApiType),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("TrafficLimit.FetchByProjectIDAndApiType"),
			),
		m,
	)
	return err
}

func (m *TrafficLimit) FetchByTrafficLimitID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("TrafficLimitID").Eq(m.TrafficLimitID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("TrafficLimit.FetchByTrafficLimitID"),
			),
		m,
	)
	return err
}

func (m *TrafficLimit) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("TrafficLimit.UpdateByIDWithFVs"),
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

func (m *TrafficLimit) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *TrafficLimit) UpdateByProjectIDAndApiTypeWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
					tbl.ColByFieldName("ApiType").Eq(m.ApiType),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("TrafficLimit.UpdateByProjectIDAndApiTypeWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByProjectIDAndApiType(db)
	}
	return nil
}

func (m *TrafficLimit) UpdateByProjectIDAndApiType(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByProjectIDAndApiTypeWithFVs(db, fvs)
}

func (m *TrafficLimit) UpdateByTrafficLimitIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("TrafficLimitID").Eq(m.TrafficLimitID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("TrafficLimit.UpdateByTrafficLimitIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByTrafficLimitID(db)
	}
	return nil
}

func (m *TrafficLimit) UpdateByTrafficLimitID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByTrafficLimitIDWithFVs(db, fvs)
}

func (m *TrafficLimit) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("TrafficLimit.Delete"),
			),
	)
	return err
}

func (m *TrafficLimit) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("TrafficLimit.DeleteByID"),
			),
	)
	return err
}

func (m *TrafficLimit) SoftDeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("TrafficLimit.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *TrafficLimit) DeleteByProjectIDAndApiType(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
						tbl.ColByFieldName("ApiType").Eq(m.ApiType),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("TrafficLimit.DeleteByProjectIDAndApiType"),
			),
	)
	return err
}

func (m *TrafficLimit) SoftDeleteByProjectIDAndApiType(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("ApiType").Eq(m.ApiType),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("TrafficLimit.SoftDeleteByProjectIDAndApiType"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *TrafficLimit) DeleteByTrafficLimitID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("TrafficLimitID").Eq(m.TrafficLimitID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("TrafficLimit.DeleteByTrafficLimitID"),
			),
	)
	return err
}

func (m *TrafficLimit) SoftDeleteByTrafficLimitID(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("TrafficLimitID").Eq(m.TrafficLimitID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("TrafficLimit.SoftDeleteByTrafficLimitID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
