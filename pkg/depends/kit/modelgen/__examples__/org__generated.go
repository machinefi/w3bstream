// This is a generated source file. DO NOT EDIT
// Source: example/org__generated.go

package example

import (
	"fmt"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var OrgTable *builder.Table

func init() {
	OrgTable = DB.Register(&Org{})
}

type OrgIterator struct {
}

func (*OrgIterator) New() interface{} {
	return &Org{}
}

func (*OrgIterator) Resolve(v interface{}) *Org {
	return v.(*Org)
}

func (*Org) TableName() string {
	return "t_org"
}

func (*Org) TableDesc() []string {
	return []string{
		"Org describe organization information",
	}
}

func (*Org) Comments() map[string]string {
	return map[string]string{
		"UserID": "User relation...",
	}
}

func (*Org) ColDesc() map[string][]string {
	return map[string][]string{
		"UserID": []string{
			"User relation...",
		},
	}
}

func (*Org) ColRel() map[string][]string {
	return map[string][]string{
		"UserID": []string{
			"User",
			"ID",
		},
	}
}

func (*Org) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Org) IndexFieldNames() []string {
	return []string{
		"ID",
	}
}

func (m *Org) ColID() *builder.Column {
	return OrgTable.ColByFieldName(m.FieldID())
}

func (*Org) FieldID() string {
	return "ID"
}

func (m *Org) ColName() *builder.Column {
	return OrgTable.ColByFieldName(m.FieldName())
}

func (*Org) FieldName() string {
	return "Name"
}

func (m *Org) ColUserID() *builder.Column {
	return OrgTable.ColByFieldName(m.FieldUserID())
}

func (*Org) FieldUserID() string {
	return "UserID"
}

func (m *Org) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *Org) Create(db sqlx.DBExecutor) error {

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Org) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Org, error) {
	var (
		tbl = db.T(m)
		lst = make([]Org, 0)
	)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Org.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Org) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Org.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Org) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Org.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Org) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ID").Eq(m.ID),
				),
				builder.Comment("Org.UpdateByIDWithFVs"),
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

func (m *Org) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Org) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Org.Delete"),
			),
	)
	return err
}

func (m *Org) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Org.DeleteByID"),
			),
	)
	return err
}
