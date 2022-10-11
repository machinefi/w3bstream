// This is a generated source file. DO NOT EDIT
// Source: models/blockchain__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/sqlx"
	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"
)

var BlockchainTable *builder.Table

func init() {
	BlockchainTable = DB.Register(&Blockchain{})
}

type BlockchainIterator struct {
}

func (BlockchainIterator) New() interface{} {
	return &Blockchain{}
}

func (BlockchainIterator) Resolve(v interface{}) *Blockchain {
	return v.(*Blockchain)
}

func (*Blockchain) TableName() string {
	return "t_blockchain"
}

func (*Blockchain) TableDesc() []string {
	return []string{
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
		"ID",
	}
}

func (m *Blockchain) ColID() *builder.Column {
	return BlockchainTable.ColByFieldName(m.FieldID())
}

func (*Blockchain) FieldID() string {
	return "ID"
}

func (m *Blockchain) ColBlockchainID() *builder.Column {
	return BlockchainTable.ColByFieldName(m.FieldBlockchainID())
}

func (*Blockchain) FieldBlockchainID() string {
	return "BlockchainID"
}

func (m *Blockchain) ColBlockchainAddress() *builder.Column {
	return BlockchainTable.ColByFieldName(m.FieldBlockchainAddress())
}

func (*Blockchain) FieldBlockchainAddress() string {
	return "BlockchainAddress"
}

func (m *Blockchain) ColContractAddress() *builder.Column {
	return BlockchainTable.ColByFieldName(m.FieldContractAddress())
}

func (*Blockchain) FieldContractAddress() string {
	return "ContractAddress"
}

func (m *Blockchain) ColBlockStart() *builder.Column {
	return BlockchainTable.ColByFieldName(m.FieldBlockStart())
}

func (*Blockchain) FieldBlockStart() string {
	return "BlockStart"
}

func (m *Blockchain) ColBlockCurrent() *builder.Column {
	return BlockchainTable.ColByFieldName(m.FieldBlockCurrent())
}

func (*Blockchain) FieldBlockCurrent() string {
	return "BlockCurrent"
}

func (m *Blockchain) ColProjectID() *builder.Column {
	return BlockchainTable.ColByFieldName(m.FieldProjectID())
}

func (*Blockchain) FieldProjectID() string {
	return "ProjectID"
}

func (m *Blockchain) ColAppletID() *builder.Column {
	return BlockchainTable.ColByFieldName(m.FieldAppletID())
}

func (*Blockchain) FieldAppletID() string {
	return "AppletID"
}

func (m *Blockchain) ColHandler() *builder.Column {
	return BlockchainTable.ColByFieldName(m.FieldHandler())
}

func (*Blockchain) FieldHandler() string {
	return "Handler"
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
