// This is a generated source file. DO NOT EDIT
// Source: models/project_operator__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var ProjectOperatorTable *builder.Table

func init() {
	ProjectOperatorTable = DB.Register(&ProjectOperator{})
}

type ProjectOperatorIterator struct {
}

func (*ProjectOperatorIterator) New() interface{} {
	return &ProjectOperator{}
}

func (*ProjectOperatorIterator) Resolve(v interface{}) *ProjectOperator {
	return v.(*ProjectOperator)
}

func (*ProjectOperator) TableName() string {
	return "t_project_operator"
}

func (*ProjectOperator) TableDesc() []string {
	return []string{
		"ProjectOperator schema for project operator relationship",
	}
}

func (*ProjectOperator) Comments() map[string]string {
	return map[string]string{}
}

func (*ProjectOperator) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*ProjectOperator) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*ProjectOperator) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *ProjectOperator) IndexFieldNames() []string {
	return []string{
		"ID",
		"ProjectID",
	}
}

func (*ProjectOperator) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_project_id": []string{
			"ProjectID",
			"DeletedAt",
		},
	}
}

func (*ProjectOperator) UniqueIndexUIProjectID() string {
	return "ui_project_id"
}

func (m *ProjectOperator) ColID() *builder.Column {
	return ProjectOperatorTable.ColByFieldName(m.FieldID())
}

func (*ProjectOperator) FieldID() string {
	return "ID"
}

func (m *ProjectOperator) ColProjectID() *builder.Column {
	return ProjectOperatorTable.ColByFieldName(m.FieldProjectID())
}

func (*ProjectOperator) FieldProjectID() string {
	return "ProjectID"
}

func (m *ProjectOperator) ColOperatorID() *builder.Column {
	return ProjectOperatorTable.ColByFieldName(m.FieldOperatorID())
}

func (*ProjectOperator) FieldOperatorID() string {
	return "OperatorID"
}

func (m *ProjectOperator) ColCreatedAt() *builder.Column {
	return ProjectOperatorTable.ColByFieldName(m.FieldCreatedAt())
}

func (*ProjectOperator) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *ProjectOperator) ColUpdatedAt() *builder.Column {
	return ProjectOperatorTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*ProjectOperator) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *ProjectOperator) ColDeletedAt() *builder.Column {
	return ProjectOperatorTable.ColByFieldName(m.FieldDeletedAt())
}

func (*ProjectOperator) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *ProjectOperator) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *ProjectOperator) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *ProjectOperator) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]ProjectOperator, error) {
	var (
		tbl = db.T(m)
		lst = make([]ProjectOperator, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("ProjectOperator.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *ProjectOperator) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("ProjectOperator.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *ProjectOperator) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("ProjectOperator.FetchByID"),
			),
		m,
	)
	return err
}

func (m *ProjectOperator) FetchByProjectID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("ProjectOperator.FetchByProjectID"),
			),
		m,
	)
	return err
}

func (m *ProjectOperator) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("ProjectOperator.UpdateByIDWithFVs"),
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

func (m *ProjectOperator) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *ProjectOperator) UpdateByProjectIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("ProjectOperator.UpdateByProjectIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByProjectID(db)
	}
	return nil
}

func (m *ProjectOperator) UpdateByProjectID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByProjectIDWithFVs(db, fvs)
}

func (m *ProjectOperator) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("ProjectOperator.Delete"),
			),
	)
	return err
}

func (m *ProjectOperator) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("ProjectOperator.DeleteByID"),
			),
	)
	return err
}

func (m *ProjectOperator) SoftDeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("ProjectOperator.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *ProjectOperator) DeleteByProjectID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("ProjectOperator.DeleteByProjectID"),
			),
	)
	return err
}

func (m *ProjectOperator) SoftDeleteByProjectID(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("ProjectOperator.SoftDeleteByProjectID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
