// This is a generated source file. DO NOT EDIT
// Source: models/operator__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var OperatorTable *builder.Table

func init() {
	OperatorTable = DB.Register(&Operator{})
}

type OperatorIterator struct {
}

func (*OperatorIterator) New() interface{} {
	return &Operator{}
}

func (*OperatorIterator) Resolve(v interface{}) *Operator {
	return v.(*Operator)
}

func (*Operator) TableName() string {
	return "t_operator"
}

func (*Operator) TableDesc() []string {
	return []string{
		"Operator schema for blockchain operate information",
	}
}

func (*Operator) Comments() map[string]string {
	return map[string]string{
		"AccountID": "AccountID  account id",
	}
}

func (*Operator) ColDesc() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"AccountID  account id",
		},
	}
}

func (*Operator) ColRel() map[string][]string {
	return map[string][]string{
		"AccountID": []string{
			"Account",
			"AccountID",
		},
	}
}

func (*Operator) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Operator) IndexFieldNames() []string {
	return []string{
		"AccountID",
		"ID",
		"Name",
		"OperatorID",
	}
}

func (*Operator) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_name": []string{
			"AccountID",
			"Name",
			"DeletedAt",
		},
		"ui_operator_id": []string{
			"OperatorID",
			"DeletedAt",
		},
	}
}

func (*Operator) UniqueIndexUIName() string {
	return "ui_name"
}

func (*Operator) UniqueIndexUIOperatorID() string {
	return "ui_operator_id"
}

func (m *Operator) ColID() *builder.Column {
	return OperatorTable.ColByFieldName(m.FieldID())
}

func (*Operator) FieldID() string {
	return "ID"
}

func (m *Operator) ColAccountID() *builder.Column {
	return OperatorTable.ColByFieldName(m.FieldAccountID())
}

func (*Operator) FieldAccountID() string {
	return "AccountID"
}

func (m *Operator) ColOperatorID() *builder.Column {
	return OperatorTable.ColByFieldName(m.FieldOperatorID())
}

func (*Operator) FieldOperatorID() string {
	return "OperatorID"
}

func (m *Operator) ColPrivateKey() *builder.Column {
	return OperatorTable.ColByFieldName(m.FieldPrivateKey())
}

func (*Operator) FieldPrivateKey() string {
	return "PrivateKey"
}

func (m *Operator) ColPaymasterKey() *builder.Column {
	return OperatorTable.ColByFieldName(m.FieldPaymasterKey())
}

func (*Operator) FieldPaymasterKey() string {
	return "PaymasterKey"
}

func (m *Operator) ColName() *builder.Column {
	return OperatorTable.ColByFieldName(m.FieldName())
}

func (*Operator) FieldName() string {
	return "Name"
}

func (m *Operator) ColType() *builder.Column {
	return OperatorTable.ColByFieldName(m.FieldType())
}

func (*Operator) FieldType() string {
	return "Type"
}

func (m *Operator) ColCreatedAt() *builder.Column {
	return OperatorTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Operator) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Operator) ColUpdatedAt() *builder.Column {
	return OperatorTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Operator) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Operator) ColDeletedAt() *builder.Column {
	return OperatorTable.ColByFieldName(m.FieldDeletedAt())
}

func (*Operator) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *Operator) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *Operator) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Operator) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Operator, error) {
	var (
		tbl = db.T(m)
		lst = make([]Operator, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Operator.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Operator) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Operator.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Operator) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Operator.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Operator) FetchByAccountIDAndName(db sqlx.DBExecutor) error {
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
				builder.Comment("Operator.FetchByAccountIDAndName"),
			),
		m,
	)
	return err
}

func (m *Operator) FetchByOperatorID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("OperatorID").Eq(m.OperatorID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Operator.FetchByOperatorID"),
			),
		m,
	)
	return err
}

func (m *Operator) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("Operator.UpdateByIDWithFVs"),
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

func (m *Operator) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Operator) UpdateByAccountIDAndNameWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("Operator.UpdateByAccountIDAndNameWithFVs"),
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

func (m *Operator) UpdateByAccountIDAndName(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByAccountIDAndNameWithFVs(db, fvs)
}

func (m *Operator) UpdateByOperatorIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("OperatorID").Eq(m.OperatorID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Operator.UpdateByOperatorIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByOperatorID(db)
	}
	return nil
}

func (m *Operator) UpdateByOperatorID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByOperatorIDWithFVs(db, fvs)
}

func (m *Operator) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Operator.Delete"),
			),
	)
	return err
}

func (m *Operator) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Operator.DeleteByID"),
			),
	)
	return err
}

func (m *Operator) SoftDeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Operator.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *Operator) DeleteByAccountIDAndName(db sqlx.DBExecutor) error {
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
				builder.Comment("Operator.DeleteByAccountIDAndName"),
			),
	)
	return err
}

func (m *Operator) SoftDeleteByAccountIDAndName(db sqlx.DBExecutor) error {
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
				builder.Comment("Operator.SoftDeleteByAccountIDAndName"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *Operator) DeleteByOperatorID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("OperatorID").Eq(m.OperatorID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Operator.DeleteByOperatorID"),
			),
	)
	return err
}

func (m *Operator) SoftDeleteByOperatorID(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("OperatorID").Eq(m.OperatorID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Operator.SoftDeleteByOperatorID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
