// This is a generated source file. DO NOT EDIT
// Source: models/resource_meta__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var ResourceMetaTable *builder.Table

func init() {
	ResourceMetaTable = DB.Register(&ResourceMeta{})
}

type ResourceMetaIterator struct {
}

func (*ResourceMetaIterator) New() interface{} {
	return &ResourceMeta{}
}

func (*ResourceMetaIterator) Resolve(v interface{}) *ResourceMeta {
	return v.(*ResourceMeta)
}

func (*ResourceMeta) TableName() string {
	return "t_resource_meta"
}

func (*ResourceMeta) TableDesc() []string {
	return []string{
		"ResourceMeta database model wasm_resource_meta",
	}
}

func (*ResourceMeta) Comments() map[string]string {
	return map[string]string{
		"AccountID": "AccountID  account id",
	}
}

func (*ResourceMeta) ColDesc() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"AccountID  account id",
		},
	}
}

func (*ResourceMeta) ColRel() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"Account",
			"AccountID",
		},
	}
}

func (*ResourceMeta) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *ResourceMeta) IndexFieldNames() []string {
	return []string{
		"AccountID",
		"AppletID",
		"ID",
		"MetaID",
		"ResourceID",
	}
}

func (*ResourceMeta) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_meta_id": []string{
			"MetaID",
			"DeletedAt",
		},
		"ui_res_acc_app": []string{
			"ResourceID",
			"AccountID",
			"AppletID",
			"DeletedAt",
		},
	}
}

func (*ResourceMeta) UniqueIndexUIMetaID() string {
	return "ui_meta_id"
}

func (*ResourceMeta) UniqueIndexUIResAccApp() string {
	return "ui_res_acc_app"
}

func (m *ResourceMeta) ColID() *builder.Column {
	return ResourceMetaTable.ColByFieldName(m.FieldID())
}

func (*ResourceMeta) FieldID() string {
	return "ID"
}

func (m *ResourceMeta) ColMetaID() *builder.Column {
	return ResourceMetaTable.ColByFieldName(m.FieldMetaID())
}

func (*ResourceMeta) FieldMetaID() string {
	return "MetaID"
}

func (m *ResourceMeta) ColResourceID() *builder.Column {
	return ResourceMetaTable.ColByFieldName(m.FieldResourceID())
}

func (*ResourceMeta) FieldResourceID() string {
	return "ResourceID"
}

func (m *ResourceMeta) ColAccountID() *builder.Column {
	return ResourceMetaTable.ColByFieldName(m.FieldAccountID())
}

func (*ResourceMeta) FieldAccountID() string {
	return "AccountID"
}

func (m *ResourceMeta) ColAppletID() *builder.Column {
	return ResourceMetaTable.ColByFieldName(m.FieldAppletID())
}

func (*ResourceMeta) FieldAppletID() string {
	return "AppletID"
}

func (m *ResourceMeta) ColFileName() *builder.Column {
	return ResourceMetaTable.ColByFieldName(m.FieldFileName())
}

func (*ResourceMeta) FieldFileName() string {
	return "FileName"
}

func (m *ResourceMeta) ColCreatedAt() *builder.Column {
	return ResourceMetaTable.ColByFieldName(m.FieldCreatedAt())
}

func (*ResourceMeta) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *ResourceMeta) ColUpdatedAt() *builder.Column {
	return ResourceMetaTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*ResourceMeta) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *ResourceMeta) ColDeletedAt() *builder.Column {
	return ResourceMetaTable.ColByFieldName(m.FieldDeletedAt())
}

func (*ResourceMeta) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *ResourceMeta) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *ResourceMeta) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *ResourceMeta) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]ResourceMeta, error) {
	var (
		tbl = db.T(m)
		lst = make([]ResourceMeta, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("ResourceMeta.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *ResourceMeta) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("ResourceMeta.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *ResourceMeta) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("ResourceMeta.FetchByID"),
			),
		m,
	)
	return err
}

func (m *ResourceMeta) FetchByMetaID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("MetaID").Eq(m.MetaID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("ResourceMeta.FetchByMetaID"),
			),
		m,
	)
	return err
}

func (m *ResourceMeta) FetchByResourceIDAndAccountIDAndAppletID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ResourceID").Eq(m.ResourceID),
						tbl.ColByFieldName("AccountID").Eq(m.AccountID),
						tbl.ColByFieldName("AppletID").Eq(m.AppletID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("ResourceMeta.FetchByResourceIDAndAccountIDAndAppletID"),
			),
		m,
	)
	return err
}

func (m *ResourceMeta) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("ResourceMeta.UpdateByIDWithFVs"),
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

func (m *ResourceMeta) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *ResourceMeta) UpdateByMetaIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("MetaID").Eq(m.MetaID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("ResourceMeta.UpdateByMetaIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByMetaID(db)
	}
	return nil
}

func (m *ResourceMeta) UpdateByMetaID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByMetaIDWithFVs(db, fvs)
}

func (m *ResourceMeta) UpdateByResourceIDAndAccountIDAndAppletIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ResourceID").Eq(m.ResourceID),
					tbl.ColByFieldName("AccountID").Eq(m.AccountID),
					tbl.ColByFieldName("AppletID").Eq(m.AppletID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("ResourceMeta.UpdateByResourceIDAndAccountIDAndAppletIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByResourceIDAndAccountIDAndAppletID(db)
	}
	return nil
}

func (m *ResourceMeta) UpdateByResourceIDAndAccountIDAndAppletID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByResourceIDAndAccountIDAndAppletIDWithFVs(db, fvs)
}

func (m *ResourceMeta) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("ResourceMeta.Delete"),
			),
	)
	return err
}

func (m *ResourceMeta) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("ResourceMeta.DeleteByID"),
			),
	)
	return err
}

func (m *ResourceMeta) SoftDeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("ResourceMeta.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *ResourceMeta) DeleteByMetaID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("MetaID").Eq(m.MetaID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("ResourceMeta.DeleteByMetaID"),
			),
	)
	return err
}

func (m *ResourceMeta) SoftDeleteByMetaID(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("MetaID").Eq(m.MetaID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("ResourceMeta.SoftDeleteByMetaID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *ResourceMeta) DeleteByResourceIDAndAccountIDAndAppletID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ResourceID").Eq(m.ResourceID),
						tbl.ColByFieldName("AccountID").Eq(m.AccountID),
						tbl.ColByFieldName("AppletID").Eq(m.AppletID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("ResourceMeta.DeleteByResourceIDAndAccountIDAndAppletID"),
			),
	)
	return err
}

func (m *ResourceMeta) SoftDeleteByResourceIDAndAccountIDAndAppletID(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("AccountID").Eq(m.AccountID),
					tbl.ColByFieldName("AppletID").Eq(m.AppletID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("ResourceMeta.SoftDeleteByResourceIDAndAccountIDAndAppletID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
