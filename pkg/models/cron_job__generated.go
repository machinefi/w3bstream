// This is a generated source file. DO NOT EDIT
// Source: models/cron_job__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var CronJobTable *builder.Table

func init() {
	CronJobTable = DB.Register(&CronJob{})
}

type CronJobIterator struct {
}

func (*CronJobIterator) New() interface{} {
	return &CronJob{}
}

func (*CronJobIterator) Resolve(v interface{}) *CronJob {
	return v.(*CronJob)
}

func (*CronJob) TableName() string {
	return "t_cron_job"
}

func (*CronJob) TableDesc() []string {
	return []string{
		"CronJob schema for cron job information",
	}
}

func (*CronJob) Comments() map[string]string {
	return map[string]string{}
}

func (*CronJob) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*CronJob) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*CronJob) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *CronJob) IndexFieldNames() []string {
	return []string{
		"CronExpressions",
		"CronJobID",
		"EventType",
		"ID",
		"ProjectID",
	}
}

func (*CronJob) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_cron": []string{
			"ProjectID",
			"CronExpressions",
			"EventType",
			"DeletedAt",
		},
		"ui_cron_job_id": []string{
			"CronJobID",
			"DeletedAt",
		},
	}
}

func (*CronJob) UniqueIndexUICron() string {
	return "ui_cron"
}

func (*CronJob) UniqueIndexUICronJobID() string {
	return "ui_cron_job_id"
}

func (m *CronJob) ColID() *builder.Column {
	return CronJobTable.ColByFieldName(m.FieldID())
}

func (*CronJob) FieldID() string {
	return "ID"
}

func (m *CronJob) ColCronJobID() *builder.Column {
	return CronJobTable.ColByFieldName(m.FieldCronJobID())
}

func (*CronJob) FieldCronJobID() string {
	return "CronJobID"
}

func (m *CronJob) ColProjectID() *builder.Column {
	return CronJobTable.ColByFieldName(m.FieldProjectID())
}

func (*CronJob) FieldProjectID() string {
	return "ProjectID"
}

func (m *CronJob) ColCronExpressions() *builder.Column {
	return CronJobTable.ColByFieldName(m.FieldCronExpressions())
}

func (*CronJob) FieldCronExpressions() string {
	return "CronExpressions"
}

func (m *CronJob) ColEventType() *builder.Column {
	return CronJobTable.ColByFieldName(m.FieldEventType())
}

func (*CronJob) FieldEventType() string {
	return "EventType"
}

func (m *CronJob) ColCreatedAt() *builder.Column {
	return CronJobTable.ColByFieldName(m.FieldCreatedAt())
}

func (*CronJob) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *CronJob) ColUpdatedAt() *builder.Column {
	return CronJobTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*CronJob) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *CronJob) ColDeletedAt() *builder.Column {
	return CronJobTable.ColByFieldName(m.FieldDeletedAt())
}

func (*CronJob) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *CronJob) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *CronJob) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *CronJob) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]CronJob, error) {
	var (
		tbl = db.T(m)
		lst = make([]CronJob, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("CronJob.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *CronJob) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("CronJob.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *CronJob) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("CronJob.FetchByID"),
			),
		m,
	)
	return err
}

func (m *CronJob) FetchByProjectIDAndCronExpressionsAndEventType(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
						tbl.ColByFieldName("CronExpressions").Eq(m.CronExpressions),
						tbl.ColByFieldName("EventType").Eq(m.EventType),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("CronJob.FetchByProjectIDAndCronExpressionsAndEventType"),
			),
		m,
	)
	return err
}

func (m *CronJob) FetchByCronJobID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("CronJobID").Eq(m.CronJobID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("CronJob.FetchByCronJobID"),
			),
		m,
	)
	return err
}

func (m *CronJob) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("CronJob.UpdateByIDWithFVs"),
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

func (m *CronJob) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *CronJob) UpdateByProjectIDAndCronExpressionsAndEventTypeWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
					tbl.ColByFieldName("CronExpressions").Eq(m.CronExpressions),
					tbl.ColByFieldName("EventType").Eq(m.EventType),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("CronJob.UpdateByProjectIDAndCronExpressionsAndEventTypeWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByProjectIDAndCronExpressionsAndEventType(db)
	}
	return nil
}

func (m *CronJob) UpdateByProjectIDAndCronExpressionsAndEventType(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByProjectIDAndCronExpressionsAndEventTypeWithFVs(db, fvs)
}

func (m *CronJob) UpdateByCronJobIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("CronJobID").Eq(m.CronJobID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("CronJob.UpdateByCronJobIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByCronJobID(db)
	}
	return nil
}

func (m *CronJob) UpdateByCronJobID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByCronJobIDWithFVs(db, fvs)
}

func (m *CronJob) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("CronJob.Delete"),
			),
	)
	return err
}

func (m *CronJob) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("CronJob.DeleteByID"),
			),
	)
	return err
}

func (m *CronJob) SoftDeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("CronJob.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *CronJob) DeleteByProjectIDAndCronExpressionsAndEventType(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
						tbl.ColByFieldName("CronExpressions").Eq(m.CronExpressions),
						tbl.ColByFieldName("EventType").Eq(m.EventType),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("CronJob.DeleteByProjectIDAndCronExpressionsAndEventType"),
			),
	)
	return err
}

func (m *CronJob) SoftDeleteByProjectIDAndCronExpressionsAndEventType(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("CronExpressions").Eq(m.CronExpressions),
					tbl.ColByFieldName("EventType").Eq(m.EventType),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("CronJob.SoftDeleteByProjectIDAndCronExpressionsAndEventType"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *CronJob) DeleteByCronJobID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("CronJobID").Eq(m.CronJobID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("CronJob.DeleteByCronJobID"),
			),
	)
	return err
}

func (m *CronJob) SoftDeleteByCronJobID(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("CronJobID").Eq(m.CronJobID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("CronJob.SoftDeleteByCronJobID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
