// This is a generated source file. DO NOT EDIT
// Source: models/account_access_key__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var AccountAccessKeyTable *builder.Table

func init() {
	AccountAccessKeyTable = DB.Register(&AccountAccessKey{})
}

type AccountAccessKeyIterator struct {
}

func (*AccountAccessKeyIterator) New() interface{} {
	return &AccountAccessKey{}
}

func (*AccountAccessKeyIterator) Resolve(v interface{}) *AccountAccessKey {
	return v.(*AccountAccessKey)
}

func (*AccountAccessKey) TableName() string {
	return "t_account_access_key"
}

func (*AccountAccessKey) TableDesc() []string {
	return []string{
		"AccountAPIKey account api access key",
	}
}

func (*AccountAccessKey) Comments() map[string]string {
	return map[string]string{
		"AccountID": "AccountID  account id",
	}
}

func (*AccountAccessKey) ColDesc() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"AccountID  account id",
		},
	}
}

func (*AccountAccessKey) ColRel() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"Account",
			"AccountID",
		},
	}
}

func (*AccountAccessKey) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *AccountAccessKey) IndexFieldNames() []string {
	return []string{
		"AccessKey",
		"AccountID",
		"ID",
		"Name",
	}
}

func (*AccountAccessKey) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_access_key": []string{
			"AccessKey",
			"DeletedAt",
		},
		"ui_account_key_name": []string{
			"AccountID",
			"Name",
			"DeletedAt",
		},
	}
}

func (*AccountAccessKey) UniqueIndexUIAccessKey() string {
	return "ui_access_key"
}

func (*AccountAccessKey) UniqueIndexUIAccountKeyName() string {
	return "ui_account_key_name"
}

func (m *AccountAccessKey) ColID() *builder.Column {
	return AccountAccessKeyTable.ColByFieldName(m.FieldID())
}

func (*AccountAccessKey) FieldID() string {
	return "ID"
}

func (m *AccountAccessKey) ColAccountID() *builder.Column {
	return AccountAccessKeyTable.ColByFieldName(m.FieldAccountID())
}

func (*AccountAccessKey) FieldAccountID() string {
	return "AccountID"
}

func (m *AccountAccessKey) ColName() *builder.Column {
	return AccountAccessKeyTable.ColByFieldName(m.FieldName())
}

func (*AccountAccessKey) FieldName() string {
	return "Name"
}

func (m *AccountAccessKey) ColAccessKey() *builder.Column {
	return AccountAccessKeyTable.ColByFieldName(m.FieldAccessKey())
}

func (*AccountAccessKey) FieldAccessKey() string {
	return "AccessKey"
}

func (m *AccountAccessKey) ColExpiredAt() *builder.Column {
	return AccountAccessKeyTable.ColByFieldName(m.FieldExpiredAt())
}

func (*AccountAccessKey) FieldExpiredAt() string {
	return "ExpiredAt"
}

func (m *AccountAccessKey) ColCreatedAt() *builder.Column {
	return AccountAccessKeyTable.ColByFieldName(m.FieldCreatedAt())
}

func (*AccountAccessKey) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *AccountAccessKey) ColUpdatedAt() *builder.Column {
	return AccountAccessKeyTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*AccountAccessKey) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *AccountAccessKey) ColDeletedAt() *builder.Column {
	return AccountAccessKeyTable.ColByFieldName(m.FieldDeletedAt())
}

func (*AccountAccessKey) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *AccountAccessKey) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *AccountAccessKey) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *AccountAccessKey) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]AccountAccessKey, error) {
	var (
		tbl = db.T(m)
		lst = make([]AccountAccessKey, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("AccountAccessKey.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *AccountAccessKey) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("AccountAccessKey.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *AccountAccessKey) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("AccountAccessKey.FetchByID"),
			),
		m,
	)
	return err
}

func (m *AccountAccessKey) FetchByAccessKey(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("AccessKey").Eq(m.AccessKey),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("AccountAccessKey.FetchByAccessKey"),
			),
		m,
	)
	return err
}

func (m *AccountAccessKey) FetchByAccountIDAndName(db sqlx.DBExecutor) error {
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
				builder.Comment("AccountAccessKey.FetchByAccountIDAndName"),
			),
		m,
	)
	return err
}

func (m *AccountAccessKey) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("AccountAccessKey.UpdateByIDWithFVs"),
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

func (m *AccountAccessKey) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *AccountAccessKey) UpdateByAccessKeyWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("AccessKey").Eq(m.AccessKey),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("AccountAccessKey.UpdateByAccessKeyWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByAccessKey(db)
	}
	return nil
}

func (m *AccountAccessKey) UpdateByAccessKey(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByAccessKeyWithFVs(db, fvs)
}

func (m *AccountAccessKey) UpdateByAccountIDAndNameWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("AccountAccessKey.UpdateByAccountIDAndNameWithFVs"),
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

func (m *AccountAccessKey) UpdateByAccountIDAndName(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByAccountIDAndNameWithFVs(db, fvs)
}

func (m *AccountAccessKey) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("AccountAccessKey.Delete"),
			),
	)
	return err
}

func (m *AccountAccessKey) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("AccountAccessKey.DeleteByID"),
			),
	)
	return err
}

func (m *AccountAccessKey) SoftDeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("AccountAccessKey.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *AccountAccessKey) DeleteByAccessKey(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("AccessKey").Eq(m.AccessKey),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("AccountAccessKey.DeleteByAccessKey"),
			),
	)
	return err
}

func (m *AccountAccessKey) SoftDeleteByAccessKey(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("AccessKey").Eq(m.AccessKey),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("AccountAccessKey.SoftDeleteByAccessKey"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *AccountAccessKey) DeleteByAccountIDAndName(db sqlx.DBExecutor) error {
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
				builder.Comment("AccountAccessKey.DeleteByAccountIDAndName"),
			),
	)
	return err
}

func (m *AccountAccessKey) SoftDeleteByAccountIDAndName(db sqlx.DBExecutor) error {
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
				builder.Comment("AccountAccessKey.SoftDeleteByAccountIDAndName"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
