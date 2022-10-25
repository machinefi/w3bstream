// This is a generated source file. DO NOT EDIT
// Source: models/chain_height__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/sqlx"
	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"
)

var ChainHeightTable *builder.Table

func init() {
	ChainHeightTable = MonitorDB.Register(&ChainHeight{})
}

type ChainHeightIterator struct {
}

func (ChainHeightIterator) New() interface{} {
	return &ChainHeight{}
}

func (ChainHeightIterator) Resolve(v interface{}) *ChainHeight {
	return v.(*ChainHeight)
}

func (*ChainHeight) TableName() string {
	return "t_chain_height"
}

func (*ChainHeight) TableDesc() []string {
	return []string{
		"ChainHeight database model chainheight",
	}
}

func (*ChainHeight) Comments() map[string]string {
	return map[string]string{}
}

func (*ChainHeight) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*ChainHeight) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*ChainHeight) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *ChainHeight) IndexFieldNames() []string {
	return []string{
		"ID",
	}
}

func (m *ChainHeight) ColID() *builder.Column {
	return ChainHeightTable.ColByFieldName(m.FieldID())
}

func (*ChainHeight) FieldID() string {
	return "ID"
}

func (m *ChainHeight) ColChainHeightID() *builder.Column {
	return ChainHeightTable.ColByFieldName(m.FieldChainHeightID())
}

func (*ChainHeight) FieldChainHeightID() string {
	return "ChainHeightID"
}

func (m *ChainHeight) ColProjectName() *builder.Column {
	return ChainHeightTable.ColByFieldName(m.FieldProjectName())
}

func (*ChainHeight) FieldProjectName() string {
	return "ProjectName"
}

func (m *ChainHeight) ColFinished() *builder.Column {
	return ChainHeightTable.ColByFieldName(m.FieldFinished())
}

func (*ChainHeight) FieldFinished() string {
	return "Finished"
}

func (m *ChainHeight) ColEventType() *builder.Column {
	return ChainHeightTable.ColByFieldName(m.FieldEventType())
}

func (*ChainHeight) FieldEventType() string {
	return "EventType"
}

func (m *ChainHeight) ColChainID() *builder.Column {
	return ChainHeightTable.ColByFieldName(m.FieldChainID())
}

func (*ChainHeight) FieldChainID() string {
	return "ChainID"
}

func (m *ChainHeight) ColHeight() *builder.Column {
	return ChainHeightTable.ColByFieldName(m.FieldHeight())
}

func (*ChainHeight) FieldHeight() string {
	return "Height"
}

func (m *ChainHeight) ColCreatedAt() *builder.Column {
	return ChainHeightTable.ColByFieldName(m.FieldCreatedAt())
}

func (*ChainHeight) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *ChainHeight) ColUpdatedAt() *builder.Column {
	return ChainHeightTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*ChainHeight) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *ChainHeight) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *ChainHeight) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *ChainHeight) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]ChainHeight, error) {
	var (
		tbl = db.T(m)
		lst = make([]ChainHeight, 0)
	)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("ChainHeight.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *ChainHeight) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("ChainHeight.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *ChainHeight) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("ChainHeight.FetchByID"),
			),
		m,
	)
	return err
}

func (m *ChainHeight) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("ChainHeight.UpdateByIDWithFVs"),
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

func (m *ChainHeight) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *ChainHeight) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("ChainHeight.Delete"),
			),
	)
	return err
}

func (m *ChainHeight) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("ChainHeight.DeleteByID"),
			),
	)
	return err
}
