// This is a generated source file. DO NOT EDIT
// Source: models/traffic_rate_limit__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var TrafficRateLimitTable *builder.Table

func init() {
	TrafficRateLimitTable = DB.Register(&TrafficRateLimit{})
}

type TrafficRateLimitIterator struct {
}

func (*TrafficRateLimitIterator) New() interface{} {
	return &TrafficRateLimit{}
}

func (*TrafficRateLimitIterator) Resolve(v interface{}) *TrafficRateLimit {
	return v.(*TrafficRateLimit)
}

func (*TrafficRateLimit) TableName() string {
	return "t_traffic_rate_limit"
}

func (*TrafficRateLimit) TableDesc() []string {
	return []string{
		"TrafficRateLimit traffic rate limit for each project",
	}
}

func (*TrafficRateLimit) Comments() map[string]string {
	return map[string]string{}
}

func (*TrafficRateLimit) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*TrafficRateLimit) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*TrafficRateLimit) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *TrafficRateLimit) IndexFieldNames() []string {
	return []string{
		"ApiType",
		"ID",
		"ProjectID",
		"RateLimitID",
	}
}

func (*TrafficRateLimit) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_prj_api_type": []string{
			"ProjectID",
			"ApiType",
			"DeletedAt",
		},
		"ui_ratelimit_id": []string{
			"RateLimitID",
			"DeletedAt",
		},
	}
}

func (*TrafficRateLimit) UniqueIndexUIPrjAPIType() string {
	return "ui_prj_api_type"
}

func (*TrafficRateLimit) UniqueIndexUIRatelimitID() string {
	return "ui_ratelimit_id"
}

func (m *TrafficRateLimit) ColID() *builder.Column {
	return TrafficRateLimitTable.ColByFieldName(m.FieldID())
}

func (*TrafficRateLimit) FieldID() string {
	return "ID"
}

func (m *TrafficRateLimit) ColRateLimitID() *builder.Column {
	return TrafficRateLimitTable.ColByFieldName(m.FieldRateLimitID())
}

func (*TrafficRateLimit) FieldRateLimitID() string {
	return "RateLimitID"
}

func (m *TrafficRateLimit) ColProjectID() *builder.Column {
	return TrafficRateLimitTable.ColByFieldName(m.FieldProjectID())
}

func (*TrafficRateLimit) FieldProjectID() string {
	return "ProjectID"
}

func (m *TrafficRateLimit) ColThreshold() *builder.Column {
	return TrafficRateLimitTable.ColByFieldName(m.FieldThreshold())
}

func (*TrafficRateLimit) FieldThreshold() string {
	return "Threshold"
}

func (m *TrafficRateLimit) ColCycleNum() *builder.Column {
	return TrafficRateLimitTable.ColByFieldName(m.FieldCycleNum())
}

func (*TrafficRateLimit) FieldCycleNum() string {
	return "CycleNum"
}

func (m *TrafficRateLimit) ColCycleUnit() *builder.Column {
	return TrafficRateLimitTable.ColByFieldName(m.FieldCycleUnit())
}

func (*TrafficRateLimit) FieldCycleUnit() string {
	return "CycleUnit"
}

func (m *TrafficRateLimit) ColApiType() *builder.Column {
	return TrafficRateLimitTable.ColByFieldName(m.FieldApiType())
}

func (*TrafficRateLimit) FieldApiType() string {
	return "ApiType"
}

func (m *TrafficRateLimit) ColCreatedAt() *builder.Column {
	return TrafficRateLimitTable.ColByFieldName(m.FieldCreatedAt())
}

func (*TrafficRateLimit) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *TrafficRateLimit) ColUpdatedAt() *builder.Column {
	return TrafficRateLimitTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*TrafficRateLimit) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *TrafficRateLimit) ColDeletedAt() *builder.Column {
	return TrafficRateLimitTable.ColByFieldName(m.FieldDeletedAt())
}

func (*TrafficRateLimit) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *TrafficRateLimit) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *TrafficRateLimit) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *TrafficRateLimit) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]TrafficRateLimit, error) {
	var (
		tbl = db.T(m)
		lst = make([]TrafficRateLimit, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("TrafficRateLimit.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *TrafficRateLimit) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("TrafficRateLimit.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *TrafficRateLimit) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("TrafficRateLimit.FetchByID"),
			),
		m,
	)
	return err
}

func (m *TrafficRateLimit) FetchByProjectIDAndApiType(db sqlx.DBExecutor) error {
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
				builder.Comment("TrafficRateLimit.FetchByProjectIDAndApiType"),
			),
		m,
	)
	return err
}

func (m *TrafficRateLimit) FetchByRateLimitID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("RateLimitID").Eq(m.RateLimitID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("TrafficRateLimit.FetchByRateLimitID"),
			),
		m,
	)
	return err
}

func (m *TrafficRateLimit) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("TrafficRateLimit.UpdateByIDWithFVs"),
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

func (m *TrafficRateLimit) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *TrafficRateLimit) UpdateByProjectIDAndApiTypeWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("TrafficRateLimit.UpdateByProjectIDAndApiTypeWithFVs"),
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

func (m *TrafficRateLimit) UpdateByProjectIDAndApiType(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByProjectIDAndApiTypeWithFVs(db, fvs)
}

func (m *TrafficRateLimit) UpdateByRateLimitIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("RateLimitID").Eq(m.RateLimitID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("TrafficRateLimit.UpdateByRateLimitIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByRateLimitID(db)
	}
	return nil
}

func (m *TrafficRateLimit) UpdateByRateLimitID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByRateLimitIDWithFVs(db, fvs)
}

func (m *TrafficRateLimit) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("TrafficRateLimit.Delete"),
			),
	)
	return err
}

func (m *TrafficRateLimit) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("TrafficRateLimit.DeleteByID"),
			),
	)
	return err
}

func (m *TrafficRateLimit) SoftDeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("TrafficRateLimit.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *TrafficRateLimit) DeleteByProjectIDAndApiType(db sqlx.DBExecutor) error {
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
				builder.Comment("TrafficRateLimit.DeleteByProjectIDAndApiType"),
			),
	)
	return err
}

func (m *TrafficRateLimit) SoftDeleteByProjectIDAndApiType(db sqlx.DBExecutor) error {
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
				builder.Comment("TrafficRateLimit.SoftDeleteByProjectIDAndApiType"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *TrafficRateLimit) DeleteByRateLimitID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("RateLimitID").Eq(m.RateLimitID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("TrafficRateLimit.DeleteByRateLimitID"),
			),
	)
	return err
}

func (m *TrafficRateLimit) SoftDeleteByRateLimitID(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("RateLimitID").Eq(m.RateLimitID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("TrafficRateLimit.SoftDeleteByRateLimitID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
