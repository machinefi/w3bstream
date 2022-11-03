// This is a generated source file. DO NOT EDIT
// Source: models/chain_tx__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var ChainTxTable *builder.Table

func init() {
	ChainTxTable = MonitorDB.Register(&ChainTx{})
}

type ChainTxIterator struct {
}

func (*ChainTxIterator) New() interface{} {
	return &ChainTx{}
}

func (*ChainTxIterator) Resolve(v interface{}) *ChainTx {
	return v.(*ChainTx)
}

func (*ChainTx) TableName() string {
	return "t_chain_tx"
}

func (*ChainTx) TableDesc() []string {
	return []string{
		"ChainTx database model chain tx",
	}
}

func (*ChainTx) Comments() map[string]string {
	return map[string]string{}
}

func (*ChainTx) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*ChainTx) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*ChainTx) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *ChainTx) IndexFieldNames() []string {
	return []string{
		"ChainID",
		"ChainTxID",
		"EventType",
		"ID",
		"ProjectName",
		"TxAddress",
		"Uniq",
	}
}

func (*ChainTx) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_chain_tx_id": []string{
			"ChainTxID",
		},
		"ui_chain_tx_uniq": []string{
			"ProjectName",
			"EventType",
			"ChainID",
			"TxAddress",
			"Uniq",
		},
	}
}

func (*ChainTx) UniqueIndexUIChainTxID() string {
	return "ui_chain_tx_id"
}

func (*ChainTx) UniqueIndexUIChainTxUniq() string {
	return "ui_chain_tx_uniq"
}

func (m *ChainTx) ColID() *builder.Column {
	return ChainTxTable.ColByFieldName(m.FieldID())
}

func (*ChainTx) FieldID() string {
	return "ID"
}

func (m *ChainTx) ColChainTxID() *builder.Column {
	return ChainTxTable.ColByFieldName(m.FieldChainTxID())
}

func (*ChainTx) FieldChainTxID() string {
	return "ChainTxID"
}

func (m *ChainTx) ColProjectName() *builder.Column {
	return ChainTxTable.ColByFieldName(m.FieldProjectName())
}

func (*ChainTx) FieldProjectName() string {
	return "ProjectName"
}

func (m *ChainTx) ColFinished() *builder.Column {
	return ChainTxTable.ColByFieldName(m.FieldFinished())
}

func (*ChainTx) FieldFinished() string {
	return "Finished"
}

func (m *ChainTx) ColUniq() *builder.Column {
	return ChainTxTable.ColByFieldName(m.FieldUniq())
}

func (*ChainTx) FieldUniq() string {
	return "Uniq"
}

func (m *ChainTx) ColEventType() *builder.Column {
	return ChainTxTable.ColByFieldName(m.FieldEventType())
}

func (*ChainTx) FieldEventType() string {
	return "EventType"
}

func (m *ChainTx) ColChainID() *builder.Column {
	return ChainTxTable.ColByFieldName(m.FieldChainID())
}

func (*ChainTx) FieldChainID() string {
	return "ChainID"
}

func (m *ChainTx) ColTxAddress() *builder.Column {
	return ChainTxTable.ColByFieldName(m.FieldTxAddress())
}

func (*ChainTx) FieldTxAddress() string {
	return "TxAddress"
}

func (m *ChainTx) ColCreatedAt() *builder.Column {
	return ChainTxTable.ColByFieldName(m.FieldCreatedAt())
}

func (*ChainTx) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *ChainTx) ColUpdatedAt() *builder.Column {
	return ChainTxTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*ChainTx) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *ChainTx) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *ChainTx) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *ChainTx) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]ChainTx, error) {
	var (
		tbl = db.T(m)
		lst = make([]ChainTx, 0)
	)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("ChainTx.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *ChainTx) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("ChainTx.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *ChainTx) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("ChainTx.FetchByID"),
			),
		m,
	)
	return err
}

func (m *ChainTx) FetchByChainTxID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ChainTxID").Eq(m.ChainTxID),
					),
				),
				builder.Comment("ChainTx.FetchByChainTxID"),
			),
		m,
	)
	return err
}

func (m *ChainTx) FetchByProjectNameAndEventTypeAndChainIDAndTxAddressAndUniq(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProjectName").Eq(m.ProjectName),
						tbl.ColByFieldName("EventType").Eq(m.EventType),
						tbl.ColByFieldName("ChainID").Eq(m.ChainID),
						tbl.ColByFieldName("TxAddress").Eq(m.TxAddress),
						tbl.ColByFieldName("Uniq").Eq(m.Uniq),
					),
				),
				builder.Comment("ChainTx.FetchByProjectNameAndEventTypeAndChainIDAndTxAddressAndUniq"),
			),
		m,
	)
	return err
}

func (m *ChainTx) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("ChainTx.UpdateByIDWithFVs"),
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

func (m *ChainTx) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *ChainTx) UpdateByChainTxIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ChainTxID").Eq(m.ChainTxID),
				),
				builder.Comment("ChainTx.UpdateByChainTxIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByChainTxID(db)
	}
	return nil
}

func (m *ChainTx) UpdateByChainTxID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByChainTxIDWithFVs(db, fvs)
}

func (m *ChainTx) UpdateByProjectNameAndEventTypeAndChainIDAndTxAddressAndUniqWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ProjectName").Eq(m.ProjectName),
					tbl.ColByFieldName("EventType").Eq(m.EventType),
					tbl.ColByFieldName("ChainID").Eq(m.ChainID),
					tbl.ColByFieldName("TxAddress").Eq(m.TxAddress),
					tbl.ColByFieldName("Uniq").Eq(m.Uniq),
				),
				builder.Comment("ChainTx.UpdateByProjectNameAndEventTypeAndChainIDAndTxAddressAndUniqWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByProjectNameAndEventTypeAndChainIDAndTxAddressAndUniq(db)
	}
	return nil
}

func (m *ChainTx) UpdateByProjectNameAndEventTypeAndChainIDAndTxAddressAndUniq(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByProjectNameAndEventTypeAndChainIDAndTxAddressAndUniqWithFVs(db, fvs)
}

func (m *ChainTx) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("ChainTx.Delete"),
			),
	)
	return err
}

func (m *ChainTx) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("ChainTx.DeleteByID"),
			),
	)
	return err
}

func (m *ChainTx) DeleteByChainTxID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ChainTxID").Eq(m.ChainTxID),
					),
				),
				builder.Comment("ChainTx.DeleteByChainTxID"),
			),
	)
	return err
}

func (m *ChainTx) DeleteByProjectNameAndEventTypeAndChainIDAndTxAddressAndUniq(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProjectName").Eq(m.ProjectName),
						tbl.ColByFieldName("EventType").Eq(m.EventType),
						tbl.ColByFieldName("ChainID").Eq(m.ChainID),
						tbl.ColByFieldName("TxAddress").Eq(m.TxAddress),
						tbl.ColByFieldName("Uniq").Eq(m.Uniq),
					),
				),
				builder.Comment("ChainTx.DeleteByProjectNameAndEventTypeAndChainIDAndTxAddressAndUniq"),
			),
	)
	return err
}
