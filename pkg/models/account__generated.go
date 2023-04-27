// This is a generated source file. DO NOT EDIT
// Source: models/account__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var AccountTable *builder.Table

func init() {
	AccountTable = DB.Register(&Account{})
}

type AccountIterator struct {
}

func (*AccountIterator) New() interface{} {
	return &Account{}
}

func (*AccountIterator) Resolve(v interface{}) *Account {
	return v.(*Account)
}

func (*Account) TableName() string {
	return "t_account"
}

func (*Account) TableDesc() []string {
	return []string{
		"Account w3bstream account",
	}
}

func (*Account) Comments() map[string]string {
	return map[string]string{
		"AccountID": "AccountID  account id",
	}
}

func (*Account) ColDesc() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"AccountID  account id",
		},
	}
}

func (*Account) ColRel() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"Account",
			"AccountID",
		},
	}
}

func (*Account) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Account) IndexFieldNames() []string {
	return []string{
		"AccountID",
		"ID",
	}
}

func (*Account) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_account_id": []string{
			"AccountID",
			"DeletedAt",
		},
	}
}

func (*Account) UniqueIndexUIAccountID() string {
	return "ui_account_id"
}

func (m *Account) ColID() *builder.Column {
	return AccountTable.ColByFieldName(m.FieldID())
}

func (*Account) FieldID() string {
	return "ID"
}

func (m *Account) ColAccountID() *builder.Column {
	return AccountTable.ColByFieldName(m.FieldAccountID())
}

func (*Account) FieldAccountID() string {
	return "AccountID"
}

func (m *Account) ColRole() *builder.Column {
	return AccountTable.ColByFieldName(m.FieldRole())
}

func (*Account) FieldRole() string {
	return "Role"
}

func (m *Account) ColState() *builder.Column {
	return AccountTable.ColByFieldName(m.FieldState())
}

func (*Account) FieldState() string {
	return "State"
}

func (m *Account) ColAvatar() *builder.Column {
	return AccountTable.ColByFieldName(m.FieldAvatar())
}

func (*Account) FieldAvatar() string {
	return "Avatar"
}

func (m *Account) ColMeta() *builder.Column {
	return AccountTable.ColByFieldName(m.FieldMeta())
}

func (*Account) FieldMeta() string {
	return "Meta"
}

func (m *Account) ColOperatorPrivateKey() *builder.Column {
	return AccountTable.ColByFieldName(m.FieldOperatorPrivateKey())
}

func (*Account) FieldOperatorPrivateKey() string {
	return "OperatorPrivateKey"
}

func (m *Account) ColCreatedAt() *builder.Column {
	return AccountTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Account) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Account) ColUpdatedAt() *builder.Column {
	return AccountTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Account) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Account) ColDeletedAt() *builder.Column {
	return AccountTable.ColByFieldName(m.FieldDeletedAt())
}

func (*Account) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *Account) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *Account) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Account) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Account, error) {
	var (
		tbl = db.T(m)
		lst = make([]Account, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Account.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Account) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Account.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Account) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Account.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Account) FetchByAccountID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("AccountID").Eq(m.AccountID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Account.FetchByAccountID"),
			),
		m,
	)
	return err
}

func (m *Account) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("Account.UpdateByIDWithFVs"),
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

func (m *Account) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Account) UpdateByAccountIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("AccountID").Eq(m.AccountID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Account.UpdateByAccountIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByAccountID(db)
	}
	return nil
}

func (m *Account) UpdateByAccountID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByAccountIDWithFVs(db, fvs)
}

func (m *Account) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Account.Delete"),
			),
	)
	return err
}

func (m *Account) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Account.DeleteByID"),
			),
	)
	return err
}

func (m *Account) SoftDeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Account.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *Account) DeleteByAccountID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("AccountID").Eq(m.AccountID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Account.DeleteByAccountID"),
			),
	)
	return err
}

func (m *Account) SoftDeleteByAccountID(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Account.SoftDeleteByAccountID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
