// This is a generated source file. DO NOT EDIT
// Source: models/chaintx__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/sqlx"
	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"
)

var ChaintxTable *builder.Table

func init() {
	ChaintxTable = DB.Register(&Chaintx{})
}

type ChaintxIterator struct {
}

func (ChaintxIterator) New() interface{} {
	return &Chaintx{}
}

func (ChaintxIterator) Resolve(v interface{}) *Chaintx {
	return v.(*Chaintx)
}

func (*Chaintx) TableName() string {
	return "t_chaintx"
}

func (*Chaintx) TableDesc() []string {
	return []string{
		"Chaintx database model chaintx",
	}
}

func (*Chaintx) Comments() map[string]string {
	return map[string]string{}
}

func (*Chaintx) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*Chaintx) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*Chaintx) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Chaintx) IndexFieldNames() []string {
	return []string{
		"ID",
	}
}

func (m *Chaintx) ColID() *builder.Column {
	return ChaintxTable.ColByFieldName(m.FieldID())
}

func (*Chaintx) FieldID() string {
	return "ID"
}

func (m *Chaintx) ColChaintxID() *builder.Column {
	return ChaintxTable.ColByFieldName(m.FieldChaintxID())
}

func (*Chaintx) FieldChaintxID() string {
	return "ChaintxID"
}

func (m *Chaintx) ColProjectName() *builder.Column {
	return ChaintxTable.ColByFieldName(m.FieldProjectName())
}

func (*Chaintx) FieldProjectName() string {
	return "ProjectName"
}

func (m *Chaintx) ColFinished() *builder.Column {
	return ChaintxTable.ColByFieldName(m.FieldFinished())
}

func (*Chaintx) FieldFinished() string {
	return "Finished"
}

func (m *Chaintx) ColEventType() *builder.Column {
	return ChaintxTable.ColByFieldName(m.FieldEventType())
}

func (*Chaintx) FieldEventType() string {
	return "EventType"
}

func (m *Chaintx) ColChainID() *builder.Column {
	return ChaintxTable.ColByFieldName(m.FieldChainID())
}

func (*Chaintx) FieldChainID() string {
	return "ChainID"
}

func (m *Chaintx) ColTxAddress() *builder.Column {
	return ChaintxTable.ColByFieldName(m.FieldTxAddress())
}

func (*Chaintx) FieldTxAddress() string {
	return "TxAddress"
}

func (m *Chaintx) ColCreatedAt() *builder.Column {
	return ChaintxTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Chaintx) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Chaintx) ColUpdatedAt() *builder.Column {
	return ChaintxTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Chaintx) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Chaintx) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *Chaintx) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Chaintx) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Chaintx, error) {
	var (
		tbl = db.T(m)
		lst = make([]Chaintx, 0)
	)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Chaintx.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Chaintx) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Chaintx.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Chaintx) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Chaintx.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Chaintx) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("Chaintx.UpdateByIDWithFVs"),
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

func (m *Chaintx) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Chaintx) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Chaintx.Delete"),
			),
	)
	return err
}

func (m *Chaintx) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Chaintx.DeleteByID"),
			),
	)
	return err
}
