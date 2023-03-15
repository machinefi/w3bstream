// This is a generated source file. DO NOT EDIT
// Source: models/account_identity__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var AccountIdentityTable *builder.Table

func init() {
	AccountIdentityTable = DB.Register(&AccountIdentity{})
}

type AccountIdentityIterator struct {
}

func (*AccountIdentityIterator) New() interface{} {
	return &AccountIdentity{}
}

func (*AccountIdentityIterator) Resolve(v interface{}) *AccountIdentity {
	return v.(*AccountIdentity)
}

func (*AccountIdentity) TableName() string {
	return "t_account_identity"
}

func (*AccountIdentity) TableDesc() []string {
	return []string{
		"AccountIdentity account identity",
	}
}

func (*AccountIdentity) Comments() map[string]string {
	return map[string]string{
		"AccountID": "AccountID  account id",
	}
}

func (*AccountIdentity) ColDesc() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"AccountID  account id",
		},
	}
}

func (*AccountIdentity) ColRel() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"Account",
			"AccountID",
		},
	}
}

func (*AccountIdentity) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (*AccountIdentity) Indexes() builder.Indexes {
	return builder.Indexes{
		"i_identity_id": []string{
			"IdentityID",
		},
		"i_source": []string{
			"Source",
		},
	}
}

func (m *AccountIdentity) IndexFieldNames() []string {
	return []string{
		"AccountID",
		"ID",
		"IdentityID",
		"Source",
		"Type",
	}
}

func (*AccountIdentity) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_account_identity": []string{
			"AccountID",
			"Type",
			"DeletedAt",
		},
		"ui_identity_id": []string{
			"Type",
			"IdentityID",
			"DeletedAt",
		},
	}
}

func (*AccountIdentity) UniqueIndexUIAccountIdentity() string {
	return "ui_account_identity"
}

func (*AccountIdentity) UniqueIndexUIIdentityID() string {
	return "ui_identity_id"
}

func (m *AccountIdentity) ColID() *builder.Column {
	return AccountIdentityTable.ColByFieldName(m.FieldID())
}

func (*AccountIdentity) FieldID() string {
	return "ID"
}

func (m *AccountIdentity) ColAccountID() *builder.Column {
	return AccountIdentityTable.ColByFieldName(m.FieldAccountID())
}

func (*AccountIdentity) FieldAccountID() string {
	return "AccountID"
}

func (m *AccountIdentity) ColType() *builder.Column {
	return AccountIdentityTable.ColByFieldName(m.FieldType())
}

func (*AccountIdentity) FieldType() string {
	return "Type"
}

func (m *AccountIdentity) ColIdentityID() *builder.Column {
	return AccountIdentityTable.ColByFieldName(m.FieldIdentityID())
}

func (*AccountIdentity) FieldIdentityID() string {
	return "IdentityID"
}

func (m *AccountIdentity) ColSource() *builder.Column {
	return AccountIdentityTable.ColByFieldName(m.FieldSource())
}

func (*AccountIdentity) FieldSource() string {
	return "Source"
}

func (m *AccountIdentity) ColMeta() *builder.Column {
	return AccountIdentityTable.ColByFieldName(m.FieldMeta())
}

func (*AccountIdentity) FieldMeta() string {
	return "Meta"
}

func (m *AccountIdentity) ColCreatedAt() *builder.Column {
	return AccountIdentityTable.ColByFieldName(m.FieldCreatedAt())
}

func (*AccountIdentity) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *AccountIdentity) ColUpdatedAt() *builder.Column {
	return AccountIdentityTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*AccountIdentity) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *AccountIdentity) ColDeletedAt() *builder.Column {
	return AccountIdentityTable.ColByFieldName(m.FieldDeletedAt())
}

func (*AccountIdentity) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *AccountIdentity) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *AccountIdentity) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *AccountIdentity) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]AccountIdentity, error) {
	var (
		tbl = db.T(m)
		lst = make([]AccountIdentity, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("AccountIdentity.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *AccountIdentity) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("AccountIdentity.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *AccountIdentity) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("AccountIdentity.FetchByID"),
			),
		m,
	)
	return err
}

func (m *AccountIdentity) FetchByAccountIDAndType(db sqlx.DBExecutor) error {
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
				builder.Comment("AccountIdentity.FetchByAccountIDAndType"),
			),
		m,
	)
	return err
}

func (m *AccountIdentity) FetchByTypeAndIdentityID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("Type").Eq(m.Type),
						tbl.ColByFieldName("IdentityID").Eq(m.IdentityID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("AccountIdentity.FetchByTypeAndIdentityID"),
			),
		m,
	)
	return err
}

func (m *AccountIdentity) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("AccountIdentity.UpdateByIDWithFVs"),
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

func (m *AccountIdentity) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *AccountIdentity) UpdateByAccountIDAndTypeWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("AccountIdentity.UpdateByAccountIDAndTypeWithFVs"),
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

func (m *AccountIdentity) UpdateByAccountIDAndType(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByAccountIDAndTypeWithFVs(db, fvs)
}

func (m *AccountIdentity) UpdateByTypeAndIdentityIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("Type").Eq(m.Type),
					tbl.ColByFieldName("IdentityID").Eq(m.IdentityID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("AccountIdentity.UpdateByTypeAndIdentityIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByTypeAndIdentityID(db)
	}
	return nil
}

func (m *AccountIdentity) UpdateByTypeAndIdentityID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByTypeAndIdentityIDWithFVs(db, fvs)
}

func (m *AccountIdentity) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("AccountIdentity.Delete"),
			),
	)
	return err
}

func (m *AccountIdentity) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("AccountIdentity.DeleteByID"),
			),
	)
	return err
}

func (m *AccountIdentity) SoftDeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("AccountIdentity.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *AccountIdentity) DeleteByAccountIDAndType(db sqlx.DBExecutor) error {
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
				builder.Comment("AccountIdentity.DeleteByAccountIDAndType"),
			),
	)
	return err
}

func (m *AccountIdentity) SoftDeleteByAccountIDAndType(db sqlx.DBExecutor) error {
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
				builder.Comment("AccountIdentity.SoftDeleteByAccountIDAndType"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *AccountIdentity) DeleteByTypeAndIdentityID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("Type").Eq(m.Type),
						tbl.ColByFieldName("IdentityID").Eq(m.IdentityID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("AccountIdentity.DeleteByTypeAndIdentityID"),
			),
	)
	return err
}

func (m *AccountIdentity) SoftDeleteByTypeAndIdentityID(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("Type").Eq(m.Type),
					tbl.ColByFieldName("IdentityID").Eq(m.IdentityID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("AccountIdentity.SoftDeleteByTypeAndIdentityID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
