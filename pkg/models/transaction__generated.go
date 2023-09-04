// This is a generated source file. DO NOT EDIT
// Source: models/transaction__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var TransactionTable *builder.Table

func init() {
	TransactionTable = DB.Register(&Transaction{})
}

type TransactionIterator struct {
}

func (*TransactionIterator) New() interface{} {
	return &Transaction{}
}

func (*TransactionIterator) Resolve(v interface{}) *Transaction {
	return v.(*Transaction)
}

func (*Transaction) TableName() string {
	return "t_transaction"
}

func (*Transaction) TableDesc() []string {
	return []string{
		"Transaction schema for blockchain transaction information",
	}
}

func (*Transaction) Comments() map[string]string {
	return map[string]string{}
}

func (*Transaction) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*Transaction) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*Transaction) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Transaction) IndexFieldNames() []string {
	return []string{
		"ID",
		"TransactionID",
	}
}

func (*Transaction) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_transaction_id": []string{
			"TransactionID",
			"DeletedAt",
		},
	}
}

func (*Transaction) UniqueIndexUITransactionID() string {
	return "ui_transaction_id"
}

func (m *Transaction) ColID() *builder.Column {
	return TransactionTable.ColByFieldName(m.FieldID())
}

func (*Transaction) FieldID() string {
	return "ID"
}

func (m *Transaction) ColTransactionID() *builder.Column {
	return TransactionTable.ColByFieldName(m.FieldTransactionID())
}

func (*Transaction) FieldTransactionID() string {
	return "TransactionID"
}

func (m *Transaction) ColProjectID() *builder.Column {
	return TransactionTable.ColByFieldName(m.FieldProjectID())
}

func (*Transaction) FieldProjectID() string {
	return "ProjectID"
}

func (m *Transaction) ColChainName() *builder.Column {
	return TransactionTable.ColByFieldName(m.FieldChainName())
}

func (*Transaction) FieldChainName() string {
	return "ChainName"
}

func (m *Transaction) ColNonce() *builder.Column {
	return TransactionTable.ColByFieldName(m.FieldNonce())
}

func (*Transaction) FieldNonce() string {
	return "Nonce"
}

func (m *Transaction) ColHash() *builder.Column {
	return TransactionTable.ColByFieldName(m.FieldHash())
}

func (*Transaction) FieldHash() string {
	return "Hash"
}

func (m *Transaction) ColSender() *builder.Column {
	return TransactionTable.ColByFieldName(m.FieldSender())
}

func (*Transaction) FieldSender() string {
	return "Sender"
}

func (m *Transaction) ColReceiver() *builder.Column {
	return TransactionTable.ColByFieldName(m.FieldReceiver())
}

func (*Transaction) FieldReceiver() string {
	return "Receiver"
}

func (m *Transaction) ColData() *builder.Column {
	return TransactionTable.ColByFieldName(m.FieldData())
}

func (*Transaction) FieldData() string {
	return "Data"
}

func (m *Transaction) ColState() *builder.Column {
	return TransactionTable.ColByFieldName(m.FieldState())
}

func (*Transaction) FieldState() string {
	return "State"
}

func (m *Transaction) ColEventType() *builder.Column {
	return TransactionTable.ColByFieldName(m.FieldEventType())
}

func (*Transaction) FieldEventType() string {
	return "EventType"
}

func (m *Transaction) ColCreatedAt() *builder.Column {
	return TransactionTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Transaction) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Transaction) ColUpdatedAt() *builder.Column {
	return TransactionTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Transaction) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Transaction) ColDeletedAt() *builder.Column {
	return TransactionTable.ColByFieldName(m.FieldDeletedAt())
}

func (*Transaction) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *Transaction) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *Transaction) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Transaction) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Transaction, error) {
	var (
		tbl = db.T(m)
		lst = make([]Transaction, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Transaction.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Transaction) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Transaction.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Transaction) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Transaction.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Transaction) FetchByTransactionID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("TransactionID").Eq(m.TransactionID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Transaction.FetchByTransactionID"),
			),
		m,
	)
	return err
}

func (m *Transaction) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("Transaction.UpdateByIDWithFVs"),
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

func (m *Transaction) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Transaction) UpdateByTransactionIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("TransactionID").Eq(m.TransactionID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Transaction.UpdateByTransactionIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByTransactionID(db)
	}
	return nil
}

func (m *Transaction) UpdateByTransactionID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByTransactionIDWithFVs(db, fvs)
}

func (m *Transaction) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Transaction.Delete"),
			),
	)
	return err
}

func (m *Transaction) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Transaction.DeleteByID"),
			),
	)
	return err
}

func (m *Transaction) SoftDeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Transaction.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *Transaction) DeleteByTransactionID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("TransactionID").Eq(m.TransactionID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Transaction.DeleteByTransactionID"),
			),
	)
	return err
}

func (m *Transaction) SoftDeleteByTransactionID(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("TransactionID").Eq(m.TransactionID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Transaction.SoftDeleteByTransactionID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
