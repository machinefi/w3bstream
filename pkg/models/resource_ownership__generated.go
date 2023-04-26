// This is a generated source file. DO NOT EDIT
// Source: models/resource_ownership__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var ResourceOwnershipTable *builder.Table

func init() {
	ResourceOwnershipTable = DB.Register(&ResourceOwnership{})
}

type ResourceOwnershipIterator struct {
}

func (*ResourceOwnershipIterator) New() interface{} {
	return &ResourceOwnership{}
}

func (*ResourceOwnershipIterator) Resolve(v interface{}) *ResourceOwnership {
	return v.(*ResourceOwnership)
}

func (*ResourceOwnership) TableName() string {
	return "t_resource_ownership"
}

func (*ResourceOwnership) TableDesc() []string {
	return []string{
		"ResourceOwnership database model resource ownership",
	}
}

func (*ResourceOwnership) Comments() map[string]string {
	return map[string]string{
		"AccountID": "AccountID  account id",
	}
}

func (*ResourceOwnership) ColDesc() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"AccountID  account id",
		},
	}
}

func (*ResourceOwnership) ColRel() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"Account",
			"AccountID",
		},
	}
}

func (*ResourceOwnership) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *ResourceOwnership) IndexFieldNames() []string {
	return []string{
		"AccountID",
		"ID",
		"ResourceID",
	}
}

func (*ResourceOwnership) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_resource_account": []string{
			"ResourceID",
			"AccountID",
		},
	}
}

func (*ResourceOwnership) UniqueIndexUIResourceAccount() string {
	return "ui_resource_account"
}

func (m *ResourceOwnership) ColID() *builder.Column {
	return ResourceOwnershipTable.ColByFieldName(m.FieldID())
}

func (*ResourceOwnership) FieldID() string {
	return "ID"
}

func (m *ResourceOwnership) ColResourceID() *builder.Column {
	return ResourceOwnershipTable.ColByFieldName(m.FieldResourceID())
}

func (*ResourceOwnership) FieldResourceID() string {
	return "ResourceID"
}

func (m *ResourceOwnership) ColAccountID() *builder.Column {
	return ResourceOwnershipTable.ColByFieldName(m.FieldAccountID())
}

func (*ResourceOwnership) FieldAccountID() string {
	return "AccountID"
}

func (m *ResourceOwnership) ColUploadedAt() *builder.Column {
	return ResourceOwnershipTable.ColByFieldName(m.FieldUploadedAt())
}

func (*ResourceOwnership) FieldUploadedAt() string {
	return "UploadedAt"
}

func (m *ResourceOwnership) ColExpireAt() *builder.Column {
	return ResourceOwnershipTable.ColByFieldName(m.FieldExpireAt())
}

func (*ResourceOwnership) FieldExpireAt() string {
	return "ExpireAt"
}

func (m *ResourceOwnership) ColFilename() *builder.Column {
	return ResourceOwnershipTable.ColByFieldName(m.FieldFilename())
}

func (*ResourceOwnership) FieldFilename() string {
	return "Filename"
}

func (m *ResourceOwnership) ColComment() *builder.Column {
	return ResourceOwnershipTable.ColByFieldName(m.FieldComment())
}

func (*ResourceOwnership) FieldComment() string {
	return "Comment"
}

func (m *ResourceOwnership) ColCreatedAt() *builder.Column {
	return ResourceOwnershipTable.ColByFieldName(m.FieldCreatedAt())
}

func (*ResourceOwnership) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *ResourceOwnership) ColUpdatedAt() *builder.Column {
	return ResourceOwnershipTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*ResourceOwnership) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *ResourceOwnership) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *ResourceOwnership) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *ResourceOwnership) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]ResourceOwnership, error) {
	var (
		tbl = db.T(m)
		lst = make([]ResourceOwnership, 0)
	)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("ResourceOwnership.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *ResourceOwnership) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("ResourceOwnership.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *ResourceOwnership) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("ResourceOwnership.FetchByID"),
			),
		m,
	)
	return err
}

func (m *ResourceOwnership) FetchByResourceIDAndAccountID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ResourceID").Eq(m.ResourceID),
						tbl.ColByFieldName("AccountID").Eq(m.AccountID),
					),
				),
				builder.Comment("ResourceOwnership.FetchByResourceIDAndAccountID"),
			),
		m,
	)
	return err
}

func (m *ResourceOwnership) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("ResourceOwnership.UpdateByIDWithFVs"),
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

func (m *ResourceOwnership) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *ResourceOwnership) UpdateByResourceIDAndAccountIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				),
				builder.Comment("ResourceOwnership.UpdateByResourceIDAndAccountIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByResourceIDAndAccountID(db)
	}
	return nil
}

func (m *ResourceOwnership) UpdateByResourceIDAndAccountID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByResourceIDAndAccountIDWithFVs(db, fvs)
}

func (m *ResourceOwnership) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("ResourceOwnership.Delete"),
			),
	)
	return err
}

func (m *ResourceOwnership) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("ResourceOwnership.DeleteByID"),
			),
	)
	return err
}

func (m *ResourceOwnership) DeleteByResourceIDAndAccountID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ResourceID").Eq(m.ResourceID),
						tbl.ColByFieldName("AccountID").Eq(m.AccountID),
					),
				),
				builder.Comment("ResourceOwnership.DeleteByResourceIDAndAccountID"),
			),
	)
	return err
}
