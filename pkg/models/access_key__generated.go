// This is a generated source file. DO NOT EDIT
// Source: models/access_key__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var AccessKeyTable *builder.Table

func init() {
	AccessKeyTable = DB.Register(&AccessKey{})
}

type AccessKeyIterator struct {
}

func (*AccessKeyIterator) New() interface{} {
	return &AccessKey{}
}

func (*AccessKeyIterator) Resolve(v interface{}) *AccessKey {
	return v.(*AccessKey)
}

func (*AccessKey) TableName() string {
	return "t_access_key"
}

func (*AccessKey) TableDesc() []string {
	return []string{
		"AccessKey api access key",
	}
}

func (*AccessKey) Comments() map[string]string {
	return map[string]string{
		"AccountID": "AccountID  account id",
	}
}

func (*AccessKey) ColDesc() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"AccountID  account id",
		},
	}
}

func (*AccessKey) ColRel() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"Account",
			"AccountID",
		},
	}
}

func (*AccessKey) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *AccessKey) IndexFieldNames() []string {
	return []string{
		"AccountID",
		"ID",
		"Name",
		"Rand",
	}
}

func (*AccessKey) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_name": []string{
			"AccountID",
			"Name",
			"DeletedAt",
		},
		"ui_rand": []string{
			"Rand",
			"DeletedAt",
		},
	}
}

func (*AccessKey) UniqueIndexUIName() string {
	return "ui_name"
}

func (*AccessKey) UniqueIndexUIRand() string {
	return "ui_rand"
}

func (m *AccessKey) ColID() *builder.Column {
	return AccessKeyTable.ColByFieldName(m.FieldID())
}

func (*AccessKey) FieldID() string {
	return "ID"
}

func (m *AccessKey) ColAccountID() *builder.Column {
	return AccessKeyTable.ColByFieldName(m.FieldAccountID())
}

func (*AccessKey) FieldAccountID() string {
	return "AccountID"
}

func (m *AccessKey) ColIdentityID() *builder.Column {
	return AccessKeyTable.ColByFieldName(m.FieldIdentityID())
}

func (*AccessKey) FieldIdentityID() string {
	return "IdentityID"
}

func (m *AccessKey) ColIdentityType() *builder.Column {
	return AccessKeyTable.ColByFieldName(m.FieldIdentityType())
}

func (*AccessKey) FieldIdentityType() string {
	return "IdentityType"
}

func (m *AccessKey) ColName() *builder.Column {
	return AccessKeyTable.ColByFieldName(m.FieldName())
}

func (*AccessKey) FieldName() string {
	return "Name"
}

func (m *AccessKey) ColRand() *builder.Column {
	return AccessKeyTable.ColByFieldName(m.FieldRand())
}

func (*AccessKey) FieldRand() string {
	return "Rand"
}

func (m *AccessKey) ColExpiredAt() *builder.Column {
	return AccessKeyTable.ColByFieldName(m.FieldExpiredAt())
}

func (*AccessKey) FieldExpiredAt() string {
	return "ExpiredAt"
}

func (m *AccessKey) ColLastUsed() *builder.Column {
	return AccessKeyTable.ColByFieldName(m.FieldLastUsed())
}

func (*AccessKey) FieldLastUsed() string {
	return "LastUsed"
}

func (m *AccessKey) ColDescription() *builder.Column {
	return AccessKeyTable.ColByFieldName(m.FieldDescription())
}

func (*AccessKey) FieldDescription() string {
	return "Description"
}

func (m *AccessKey) ColCreatedAt() *builder.Column {
	return AccessKeyTable.ColByFieldName(m.FieldCreatedAt())
}

func (*AccessKey) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *AccessKey) ColUpdatedAt() *builder.Column {
	return AccessKeyTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*AccessKey) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *AccessKey) ColDeletedAt() *builder.Column {
	return AccessKeyTable.ColByFieldName(m.FieldDeletedAt())
}

func (*AccessKey) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *AccessKey) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *AccessKey) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *AccessKey) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]AccessKey, error) {
	var (
		tbl = db.T(m)
		lst = make([]AccessKey, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("AccessKey.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *AccessKey) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("AccessKey.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *AccessKey) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("AccessKey.FetchByID"),
			),
		m,
	)
	return err
}

func (m *AccessKey) FetchByAccountIDAndName(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("AccountID").Eq(m.AccountID),
						tbl.ColByFieldName("Name").Eq(m.Name),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("AccessKey.FetchByAccountIDAndName"),
			),
		m,
	)
	return err
}

func (m *AccessKey) FetchByRand(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("Rand").Eq(m.Rand),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("AccessKey.FetchByRand"),
			),
		m,
	)
	return err
}

func (m *AccessKey) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("AccessKey.UpdateByIDWithFVs"),
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

func (m *AccessKey) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *AccessKey) UpdateByAccountIDAndNameWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("AccountID").Eq(m.AccountID),
					tbl.ColByFieldName("Name").Eq(m.Name),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("AccessKey.UpdateByAccountIDAndNameWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByAccountIDAndName(db)
	}
	return nil
}

func (m *AccessKey) UpdateByAccountIDAndName(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByAccountIDAndNameWithFVs(db, fvs)
}

func (m *AccessKey) UpdateByRandWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("Rand").Eq(m.Rand),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("AccessKey.UpdateByRandWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByRand(db)
	}
	return nil
}

func (m *AccessKey) UpdateByRand(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByRandWithFVs(db, fvs)
}

func (m *AccessKey) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("AccessKey.Delete"),
			),
	)
	return err
}

func (m *AccessKey) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("AccessKey.DeleteByID"),
			),
	)
	return err
}

func (m *AccessKey) SoftDeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("AccessKey.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *AccessKey) DeleteByAccountIDAndName(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("AccountID").Eq(m.AccountID),
						tbl.ColByFieldName("Name").Eq(m.Name),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("AccessKey.DeleteByAccountIDAndName"),
			),
	)
	return err
}

func (m *AccessKey) SoftDeleteByAccountIDAndName(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("AccountID").Eq(m.AccountID),
					tbl.ColByFieldName("Name").Eq(m.Name),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("AccessKey.SoftDeleteByAccountIDAndName"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *AccessKey) DeleteByRand(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("Rand").Eq(m.Rand),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("AccessKey.DeleteByRand"),
			),
	)
	return err
}

func (m *AccessKey) SoftDeleteByRand(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("Rand").Eq(m.Rand),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("AccessKey.SoftDeleteByRand"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
