// This is a generated source file. DO NOT EDIT
// Source: models/contractlog__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/sqlx"
	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"
)

var ContractlogTable *builder.Table

func init() {
	ContractlogTable = DB.Register(&Contractlog{})
}

type ContractlogIterator struct {
}

func (ContractlogIterator) New() interface{} {
	return &Contractlog{}
}

func (ContractlogIterator) Resolve(v interface{}) *Contractlog {
	return v.(*Contractlog)
}

func (*Contractlog) TableName() string {
	return "t_contractlog"
}

func (*Contractlog) TableDesc() []string {
	return []string{
		"Contractlog database model contractlog",
	}
}

func (*Contractlog) Comments() map[string]string {
	return map[string]string{}
}

func (*Contractlog) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*Contractlog) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*Contractlog) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Contractlog) IndexFieldNames() []string {
	return []string{
		"ID",
	}
}

func (m *Contractlog) ColID() *builder.Column {
	return ContractlogTable.ColByFieldName(m.FieldID())
}

func (*Contractlog) FieldID() string {
	return "ID"
}

func (m *Contractlog) ColContractlogID() *builder.Column {
	return ContractlogTable.ColByFieldName(m.FieldContractlogID())
}

func (*Contractlog) FieldContractlogID() string {
	return "ContractlogID"
}

func (m *Contractlog) ColProjectName() *builder.Column {
	return ContractlogTable.ColByFieldName(m.FieldProjectName())
}

func (*Contractlog) FieldProjectName() string {
	return "ProjectName"
}

func (m *Contractlog) ColEventType() *builder.Column {
	return ContractlogTable.ColByFieldName(m.FieldEventType())
}

func (*Contractlog) FieldEventType() string {
	return "EventType"
}

func (m *Contractlog) ColChainID() *builder.Column {
	return ContractlogTable.ColByFieldName(m.FieldChainID())
}

func (*Contractlog) FieldChainID() string {
	return "ChainID"
}

func (m *Contractlog) ColContractAddress() *builder.Column {
	return ContractlogTable.ColByFieldName(m.FieldContractAddress())
}

func (*Contractlog) FieldContractAddress() string {
	return "ContractAddress"
}

func (m *Contractlog) ColBlockStart() *builder.Column {
	return ContractlogTable.ColByFieldName(m.FieldBlockStart())
}

func (*Contractlog) FieldBlockStart() string {
	return "BlockStart"
}

func (m *Contractlog) ColBlockCurrent() *builder.Column {
	return ContractlogTable.ColByFieldName(m.FieldBlockCurrent())
}

func (*Contractlog) FieldBlockCurrent() string {
	return "BlockCurrent"
}

func (m *Contractlog) ColBlockEnd() *builder.Column {
	return ContractlogTable.ColByFieldName(m.FieldBlockEnd())
}

func (*Contractlog) FieldBlockEnd() string {
	return "BlockEnd"
}

func (m *Contractlog) ColTopic0() *builder.Column {
	return ContractlogTable.ColByFieldName(m.FieldTopic0())
}

func (*Contractlog) FieldTopic0() string {
	return "Topic0"
}

func (m *Contractlog) ColTopic1() *builder.Column {
	return ContractlogTable.ColByFieldName(m.FieldTopic1())
}

func (*Contractlog) FieldTopic1() string {
	return "Topic1"
}

func (m *Contractlog) ColTopic2() *builder.Column {
	return ContractlogTable.ColByFieldName(m.FieldTopic2())
}

func (*Contractlog) FieldTopic2() string {
	return "Topic2"
}

func (m *Contractlog) ColTopic3() *builder.Column {
	return ContractlogTable.ColByFieldName(m.FieldTopic3())
}

func (*Contractlog) FieldTopic3() string {
	return "Topic3"
}

func (m *Contractlog) ColCreatedAt() *builder.Column {
	return ContractlogTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Contractlog) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Contractlog) ColUpdatedAt() *builder.Column {
	return ContractlogTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Contractlog) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Contractlog) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *Contractlog) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Contractlog) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Contractlog, error) {
	var (
		tbl = db.T(m)
		lst = make([]Contractlog, 0)
	)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Contractlog.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Contractlog) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Contractlog.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Contractlog) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Contractlog.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Contractlog) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("Contractlog.UpdateByIDWithFVs"),
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

func (m *Contractlog) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Contractlog) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Contractlog.Delete"),
			),
	)
	return err
}

func (m *Contractlog) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Contractlog.DeleteByID"),
			),
	)
	return err
}
