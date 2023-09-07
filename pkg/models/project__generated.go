// This is a generated source file. DO NOT EDIT
// Source: models/project__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var ProjectTable *builder.Table

func init() {
	ProjectTable = DB.Register(&Project{})
}

type ProjectIterator struct {
}

func (*ProjectIterator) New() interface{} {
	return &Project{}
}

func (*ProjectIterator) Resolve(v interface{}) *Project {
	return v.(*Project)
}

func (*Project) TableName() string {
	return "t_project"
}

func (*Project) TableDesc() []string {
	return []string{
		"Project schema for project information",
	}
}

func (*Project) Comments() map[string]string {
	return map[string]string{
		"AccountID": "AccountID  account id",
		"Name":      "Name project name",
		"Proto":     "Proto project protocol for event publisher",
		"Public":    "Public is true, project receive event from anonymous publisher",
		"Version":   "Version project version",
	}
}

func (*Project) ColDesc() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"AccountID  account id",
		},
		"Name": []string{
			"Name project name",
		},
		"Proto": []string{
			"Proto project protocol for event publisher",
		},
		"Public": []string{
			"Public is true, project receive event from anonymous publisher",
		},
		"Version": []string{
			"Version project version",
		},
	}
}

func (*Project) ColRel() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"Account",
			"AccountID",
		},
	}
}

func (*Project) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Project) IndexFieldNames() []string {
	return []string{
		"ID",
		"Name",
		"ProjectID",
	}
}

func (*Project) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_name": []string{
			"Name",
			"DeletedAt",
		},
		"ui_project_id": []string{
			"ProjectID",
			"DeletedAt",
		},
	}
}

func (*Project) UniqueIndexUIName() string {
	return "ui_name"
}

func (*Project) UniqueIndexUIProjectID() string {
	return "ui_project_id"
}

func (m *Project) ColID() *builder.Column {
	return ProjectTable.ColByFieldName(m.FieldID())
}

func (*Project) FieldID() string {
	return "ID"
}

func (m *Project) ColProjectID() *builder.Column {
	return ProjectTable.ColByFieldName(m.FieldProjectID())
}

func (*Project) FieldProjectID() string {
	return "ProjectID"
}

func (m *Project) ColAccountID() *builder.Column {
	return ProjectTable.ColByFieldName(m.FieldAccountID())
}

func (*Project) FieldAccountID() string {
	return "AccountID"
}

func (m *Project) ColName() *builder.Column {
	return ProjectTable.ColByFieldName(m.FieldName())
}

func (*Project) FieldName() string {
	return "Name"
}

func (m *Project) ColPublic() *builder.Column {
	return ProjectTable.ColByFieldName(m.FieldPublic())
}

func (*Project) FieldPublic() string {
	return "Public"
}

func (m *Project) ColVersion() *builder.Column {
	return ProjectTable.ColByFieldName(m.FieldVersion())
}

func (*Project) FieldVersion() string {
	return "Version"
}

func (m *Project) ColProto() *builder.Column {
	return ProjectTable.ColByFieldName(m.FieldProto())
}

func (*Project) FieldProto() string {
	return "Proto"
}

func (m *Project) ColDescription() *builder.Column {
	return ProjectTable.ColByFieldName(m.FieldDescription())
}

func (*Project) FieldDescription() string {
	return "Description"
}

func (m *Project) ColCreatedAt() *builder.Column {
	return ProjectTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Project) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Project) ColUpdatedAt() *builder.Column {
	return ProjectTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Project) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Project) ColDeletedAt() *builder.Column {
	return ProjectTable.ColByFieldName(m.FieldDeletedAt())
}

func (*Project) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *Project) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *Project) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Project) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Project, error) {
	var (
		tbl = db.T(m)
		lst = make([]Project, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Project.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Project) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Project.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Project) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Project.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Project) FetchByName(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("Name").Eq(m.Name),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Project.FetchByName"),
			),
		m,
	)
	return err
}

func (m *Project) FetchByProjectID(db sqlx.DBExecutor) error {
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
				builder.Comment("Project.FetchByProjectID"),
			),
		m,
	)
	return err
}

func (m *Project) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("Project.UpdateByIDWithFVs"),
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

func (m *Project) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Project) UpdateByNameWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("Name").Eq(m.Name),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Project.UpdateByNameWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByName(db)
	}
	return nil
}

func (m *Project) UpdateByName(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByNameWithFVs(db, fvs)
}

func (m *Project) UpdateByProjectIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("Project.UpdateByProjectIDWithFVs"),
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

func (m *Project) UpdateByProjectID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByProjectIDWithFVs(db, fvs)
}

func (m *Project) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Project.Delete"),
			),
	)
	return err
}

func (m *Project) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Project.DeleteByID"),
			),
	)
	return err
}

func (m *Project) SoftDeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Project.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *Project) DeleteByName(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("Name").Eq(m.Name),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Project.DeleteByName"),
			),
	)
	return err
}

func (m *Project) SoftDeleteByName(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("Name").Eq(m.Name),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Project.SoftDeleteByName"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *Project) DeleteByProjectID(db sqlx.DBExecutor) error {
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
				builder.Comment("Project.DeleteByProjectID"),
			),
	)
	return err
}

func (m *Project) SoftDeleteByProjectID(db sqlx.DBExecutor) error {
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
				builder.Comment("Project.SoftDeleteByProjectID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
