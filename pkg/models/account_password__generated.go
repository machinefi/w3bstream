// This is a generated source file. DO NOT EDIT
// Source: models/account_password__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var AccountPasswordTable *builder.Table

func init() {
	AccountPasswordTable = DB.Register(&AccountPassword{})
}

type AccountPasswordIterator struct {
}

func (*AccountPasswordIterator) New() interface{} {
	return &AccountPassword{}
}

func (*AccountPasswordIterator) Resolve(v interface{}) *AccountPassword {
	return v.(*AccountPassword)
}

func (*AccountPassword) TableName() string {
	return "t_account_password"
}

func (*AccountPassword) TableDesc() []string {
	return []string{
		"AccountPassword account password",
	}
}

func (*AccountPassword) Comments() map[string]string {
	return map[string]string{
		"AccountID": "AccountID  account id",
		"Password":  "Password md5(md5(${account_id}-${password}))",
		"Remark":    "Remark",
		"Scope":     "Scope comma separated",
		"Type":      "Type password type",
	}
}

func (*AccountPassword) ColDesc() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"AccountID  account id",
		},
		"Password": []string{
			"Password md5(md5(${account_id}-${password}))",
		},
		"Remark": []string{
			"Remark",
		},
		"Scope": []string{
			"Scope comma separated",
		},
		"Type": []string{
			"Type password type",
		},
	}
}

func (*AccountPassword) ColRel() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"Account",
			"AccountID",
		},
	}
}

func (*AccountPassword) PrimaryKey() []string {
	return []string{
		"PasswordID",
	}
}

func (m *AccountPassword) IndexFieldNames() []string {
	return []string{
		"AccountID",
		"PasswordID",
		"Type",
	}
}

func (*AccountPassword) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_account_password": []string{
			"AccountID",
			"Type",
			"DeletedAt",
		},
	}
}

func (*AccountPassword) UniqueIndexUIAccountPassword() string {
	return "ui_account_password"
}

func (m *AccountPassword) ColAccountID() *builder.Column {
	return AccountPasswordTable.ColByFieldName(m.FieldAccountID())
}

func (*AccountPassword) FieldAccountID() string {
	return "AccountID"
}

func (m *AccountPassword) ColPasswordID() *builder.Column {
	return AccountPasswordTable.ColByFieldName(m.FieldPasswordID())
}

func (*AccountPassword) FieldPasswordID() string {
	return "PasswordID"
}

func (m *AccountPassword) ColType() *builder.Column {
	return AccountPasswordTable.ColByFieldName(m.FieldType())
}

func (*AccountPassword) FieldType() string {
	return "Type"
}

func (m *AccountPassword) ColPassword() *builder.Column {
	return AccountPasswordTable.ColByFieldName(m.FieldPassword())
}

func (*AccountPassword) FieldPassword() string {
	return "Password"
}

func (m *AccountPassword) ColScope() *builder.Column {
	return AccountPasswordTable.ColByFieldName(m.FieldScope())
}

func (*AccountPassword) FieldScope() string {
	return "Scope"
}

func (m *AccountPassword) ColRemark() *builder.Column {
	return AccountPasswordTable.ColByFieldName(m.FieldRemark())
}

func (*AccountPassword) FieldRemark() string {
	return "Remark"
}

func (m *AccountPassword) ColCreatedAt() *builder.Column {
	return AccountPasswordTable.ColByFieldName(m.FieldCreatedAt())
}

func (*AccountPassword) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *AccountPassword) ColUpdatedAt() *builder.Column {
	return AccountPasswordTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*AccountPassword) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *AccountPassword) ColDeletedAt() *builder.Column {
	return AccountPasswordTable.ColByFieldName(m.FieldDeletedAt())
}

func (*AccountPassword) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *AccountPassword) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *AccountPassword) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *AccountPassword) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]AccountPassword, error) {
	var (
		tbl = db.T(m)
		lst = make([]AccountPassword, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("AccountPassword.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *AccountPassword) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("AccountPassword.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *AccountPassword) FetchByPasswordID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("PasswordID").Eq(m.PasswordID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("AccountPassword.FetchByPasswordID"),
			),
		m,
	)
	return err
}

func (m *AccountPassword) FetchByAccountIDAndType(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("AccountID").Eq(m.AccountID),
						tbl.ColByFieldName("Type").Eq(m.Type),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("AccountPassword.FetchByAccountIDAndType"),
			),
		m,
	)
	return err
}

func (m *AccountPassword) UpdateByPasswordIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("PasswordID").Eq(m.PasswordID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("AccountPassword.UpdateByPasswordIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByPasswordID(db)
	}
	return nil
}

func (m *AccountPassword) UpdateByPasswordID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByPasswordIDWithFVs(db, fvs)
}

func (m *AccountPassword) UpdateByAccountIDAndTypeWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("AccountID").Eq(m.AccountID),
					tbl.ColByFieldName("Type").Eq(m.Type),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("AccountPassword.UpdateByAccountIDAndTypeWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByAccountIDAndType(db)
	}
	return nil
}

func (m *AccountPassword) UpdateByAccountIDAndType(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByAccountIDAndTypeWithFVs(db, fvs)
}

func (m *AccountPassword) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("AccountPassword.Delete"),
			),
	)
	return err
}

func (m *AccountPassword) DeleteByPasswordID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("PasswordID").Eq(m.PasswordID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("AccountPassword.DeleteByPasswordID"),
			),
	)
	return err
}

func (m *AccountPassword) SoftDeleteByPasswordID(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("PasswordID").Eq(m.PasswordID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("AccountPassword.SoftDeleteByPasswordID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *AccountPassword) DeleteByAccountIDAndType(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("AccountID").Eq(m.AccountID),
						tbl.ColByFieldName("Type").Eq(m.Type),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("AccountPassword.DeleteByAccountIDAndType"),
			),
	)
	return err
}

func (m *AccountPassword) SoftDeleteByAccountIDAndType(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("Type").Eq(m.Type),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("AccountPassword.SoftDeleteByAccountIDAndType"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
