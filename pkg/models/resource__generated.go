// This is a generated source file. DO NOT EDIT
// Source: models/resource__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var ResourceTable *builder.Table

func init() {
	ResourceTable = DB.Register(&Resource{})
}

type ResourceIterator struct {
}

func (*ResourceIterator) New() interface{} {
	return &Resource{}
}

func (*ResourceIterator) Resolve(v interface{}) *Resource {
	return v.(*Resource)
}

func (*Resource) TableName() string {
	return "t_resource"
}

func (*Resource) TableDesc() []string {
	return []string{
		"Resource database model wasm_resource",
	}
}

func (*Resource) Comments() map[string]string {
	return map[string]string{
		"Path": "Path <=> md5",
	}
}

func (*Resource) ColDesc() map[string][]string {
	return map[string][]string{
		"Path": []string{
			"Path <=> md5",
		},
	}
}

func (*Resource) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*Resource) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Resource) IndexFieldNames() []string {
	return []string{
		"ID",
		"Path",
		"ResourceID",
	}
}

func (*Resource) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_path": []string{
			"Path",
			"DeletedAt",
		},
		"ui_resource_id": []string{
			"ResourceID",
			"DeletedAt",
		},
	}
}

func (*Resource) UniqueIndexUIPath() string {
	return "ui_path"
}

func (*Resource) UniqueIndexUIResourceID() string {
	return "ui_resource_id"
}

func (m *Resource) ColID() *builder.Column {
	return ResourceTable.ColByFieldName(m.FieldID())
}

func (*Resource) FieldID() string {
	return "ID"
}

func (m *Resource) ColResourceID() *builder.Column {
	return ResourceTable.ColByFieldName(m.FieldResourceID())
}

func (*Resource) FieldResourceID() string {
	return "ResourceID"
}

func (m *Resource) ColPath() *builder.Column {
	return ResourceTable.ColByFieldName(m.FieldPath())
}

func (*Resource) FieldPath() string {
	return "Path"
}

func (m *Resource) ColCreatedAt() *builder.Column {
	return ResourceTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Resource) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Resource) ColUpdatedAt() *builder.Column {
	return ResourceTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Resource) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Resource) ColDeletedAt() *builder.Column {
	return ResourceTable.ColByFieldName(m.FieldDeletedAt())
}

func (*Resource) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *Resource) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *Resource) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Resource) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Resource, error) {
	var (
		tbl = db.T(m)
		lst = make([]Resource, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Resource.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Resource) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Resource.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Resource) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Resource.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Resource) FetchByPath(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("Path").Eq(m.Path),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Resource.FetchByPath"),
			),
		m,
	)
	return err
}

func (m *Resource) FetchByResourceID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ResourceID").Eq(m.ResourceID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Resource.FetchByResourceID"),
			),
		m,
	)
	return err
}

func (m *Resource) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("Resource.UpdateByIDWithFVs"),
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

func (m *Resource) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Resource) UpdateByPathWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("Path").Eq(m.Path),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Resource.UpdateByPathWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByPath(db)
	}
	return nil
}

func (m *Resource) UpdateByPath(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByPathWithFVs(db, fvs)
}

func (m *Resource) UpdateByResourceIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ResourceID").Eq(m.ResourceID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Resource.UpdateByResourceIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByResourceID(db)
	}
	return nil
}

func (m *Resource) UpdateByResourceID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByResourceIDWithFVs(db, fvs)
}

func (m *Resource) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Resource.Delete"),
			),
	)
	return err
}

func (m *Resource) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Resource.DeleteByID"),
			),
	)
	return err
}

func (m *Resource) SoftDeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Resource.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *Resource) DeleteByPath(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("Path").Eq(m.Path),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Resource.DeleteByPath"),
			),
	)
	return err
}

func (m *Resource) SoftDeleteByPath(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("Path").Eq(m.Path),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Resource.SoftDeleteByPath"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *Resource) DeleteByResourceID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ResourceID").Eq(m.ResourceID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Resource.DeleteByResourceID"),
			),
	)
	return err
}

func (m *Resource) SoftDeleteByResourceID(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("ResourceID").Eq(m.ResourceID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Resource.SoftDeleteByResourceID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
