// This is a generated source file. DO NOT EDIT
// Source: models/blockchain__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var BlockchainTable *builder.Table

func init() {
	BlockchainTable = MonitorDB.Register(&Blockchain{})
}

type BlockchainIterator struct {
}

func (*BlockchainIterator) New() interface{} {
	return &Blockchain{}
}

func (*BlockchainIterator) Resolve(v interface{}) *Blockchain {
	return v.(*Blockchain)
}

func (*Blockchain) TableName() string {
	return "t_blockchain"
}

func (*Blockchain) TableDesc() []string {
	return []string{
		"TODO will delete after caller migrate to http api",
		"Blockchain database model blockchain",
	}
}

func (*Blockchain) Comments() map[string]string {
	return map[string]string{}
}

func (*Blockchain) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*Blockchain) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*Blockchain) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Blockchain) IndexFieldNames() []string {
	return []string{
		"ChainID",
		"ID",
	}
}

func (*Blockchain) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_chain_id": []string{
			"ChainID",
		},
	}
}

func (*Blockchain) UniqueIndexUIChainID() string {
	return "ui_chain_id"
}

func (m *Blockchain) ColID() *builder.Column {
	return BlockchainTable.ColByFieldName(m.FieldID())
}

func (*Blockchain) FieldID() string {
	return "ID"
}

func (m *Blockchain) ColChainID() *builder.Column {
	return BlockchainTable.ColByFieldName(m.FieldChainID())
}

func (*Blockchain) FieldChainID() string {
	return "ChainID"
}

func (m *Blockchain) ColAddress() *builder.Column {
	return BlockchainTable.ColByFieldName(m.FieldAddress())
}

func (*Blockchain) FieldAddress() string {
	return "Address"
}

func (m *Blockchain) ColCreatedAt() *builder.Column {
	return BlockchainTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Blockchain) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Blockchain) ColUpdatedAt() *builder.Column {
	return BlockchainTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Blockchain) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Blockchain) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *Blockchain) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Blockchain) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Blockchain, error) {
	var (
		tbl = db.T(m)
		lst = make([]Blockchain, 0)
	)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Blockchain.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Blockchain) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Blockchain.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Blockchain) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Blockchain.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Blockchain) FetchByChainID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ChainID").Eq(m.ChainID),
					),
				),
				builder.Comment("Blockchain.FetchByChainID"),
			),
		m,
	)
	return err
}

func (m *Blockchain) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("Blockchain.UpdateByIDWithFVs"),
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

func (m *Blockchain) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Blockchain) UpdateByChainIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ChainID").Eq(m.ChainID),
				),
				builder.Comment("Blockchain.UpdateByChainIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByChainID(db)
	}
	return nil
}

func (m *Blockchain) UpdateByChainID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByChainIDWithFVs(db, fvs)
}

func (m *Blockchain) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Blockchain.Delete"),
			),
	)
	return err
}

func (m *Blockchain) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Blockchain.DeleteByID"),
			),
	)
	return err
}

func (m *Blockchain) DeleteByChainID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ChainID").Eq(m.ChainID),
					),
				),
				builder.Comment("Blockchain.DeleteByChainID"),
			),
	)
	return err
}
