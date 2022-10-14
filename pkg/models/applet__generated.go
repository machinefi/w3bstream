// This is a generated source file. DO NOT EDIT
// Source: models/applet__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/sqlx"
	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"
)

var AppletTable *builder.Table

func init() {
	AppletTable = DB.Register(&Applet{})
}

type AppletIterator struct {
}

func (AppletIterator) New() interface{} {
	return &Applet{}
}

func (AppletIterator) Resolve(v interface{}) *Applet {
	return v.(*Applet)
}

func (*Applet) TableName() string {
	return "t_applet"
}

func (*Applet) TableDesc() []string {
	return []string{
		"Applet database model applet",
	}
}

func (*Applet) Comments() map[string]string {
	return map[string]string{}
}

func (*Applet) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*Applet) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*Applet) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Applet) IndexFieldNames() []string {
	return []string{
		"AppletID",
		"ID",
		"Md5",
		"Name",
	}
}

func (*Applet) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_applet_id": []string{
			"AppletID",
		},
		"ui_md5": []string{
			"Md5",
		},
		"ui_name": []string{
			"Name",
		},
	}
}

func (*Applet) UniqueIndexUIAppletID() string {
	return "ui_applet_id"
}

func (*Applet) UniqueIndexUIMd5() string {
	return "ui_md5"
}

func (*Applet) UniqueIndexUIName() string {
	return "ui_name"
}

func (m *Applet) ColID() *builder.Column {
	return AppletTable.ColByFieldName(m.FieldID())
}

func (*Applet) FieldID() string {
	return "ID"
}

func (m *Applet) ColProjectID() *builder.Column {
	return AppletTable.ColByFieldName(m.FieldProjectID())
}

func (*Applet) FieldProjectID() string {
	return "ProjectID"
}

func (m *Applet) ColAppletID() *builder.Column {
	return AppletTable.ColByFieldName(m.FieldAppletID())
}

func (*Applet) FieldAppletID() string {
	return "AppletID"
}

func (m *Applet) ColName() *builder.Column {
	return AppletTable.ColByFieldName(m.FieldName())
}

func (*Applet) FieldName() string {
	return "Name"
}

func (m *Applet) ColPath() *builder.Column {
	return AppletTable.ColByFieldName(m.FieldPath())
}

func (*Applet) FieldPath() string {
	return "Path"
}

func (m *Applet) ColMd5() *builder.Column {
	return AppletTable.ColByFieldName(m.FieldMd5())
}

func (*Applet) FieldMd5() string {
	return "Md5"
}

func (m *Applet) ColCreatedAt() *builder.Column {
	return AppletTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Applet) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Applet) ColUpdatedAt() *builder.Column {
	return AppletTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Applet) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Applet) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *Applet) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Applet) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Applet, error) {
	var (
		tbl = db.T(m)
		lst = make([]Applet, 0)
	)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Applet.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Applet) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Applet.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Applet) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Applet.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Applet) FetchByName(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("Name").Eq(m.Name),
					),
				),
				builder.Comment("Applet.FetchByName"),
			),
		m,
	)
	return err
}

func (m *Applet) FetchByMd5(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("Md5").Eq(m.Md5),
					),
				),
				builder.Comment("Applet.FetchByMd5"),
			),
		m,
	)
	return err
}

func (m *Applet) FetchByAppletID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("AppletID").Eq(m.AppletID),
					),
				),
				builder.Comment("Applet.FetchByAppletID"),
			),
		m,
	)
	return err
}

func (m *Applet) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("Applet.UpdateByIDWithFVs"),
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

func (m *Applet) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Applet) UpdateByNameWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("Name").Eq(m.Name),
				),
				builder.Comment("Applet.UpdateByNameWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByName(db)
	}
	return nil
}

func (m *Applet) UpdateByName(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByNameWithFVs(db, fvs)
}

func (m *Applet) UpdateByMd5WithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("Md5").Eq(m.Md5),
				),
				builder.Comment("Applet.UpdateByMd5WithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByMd5(db)
	}
	return nil
}

func (m *Applet) UpdateByMd5(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByMd5WithFVs(db, fvs)
}

func (m *Applet) UpdateByAppletIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("AppletID").Eq(m.AppletID),
				),
				builder.Comment("Applet.UpdateByAppletIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByAppletID(db)
	}
	return nil
}

func (m *Applet) UpdateByAppletID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByAppletIDWithFVs(db, fvs)
}

func (m *Applet) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Applet.Delete"),
			),
	)
	return err
}

func (m *Applet) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Applet.DeleteByID"),
			),
	)
	return err
}

func (m *Applet) DeleteByName(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("Name").Eq(m.Name),
					),
				),
				builder.Comment("Applet.DeleteByName"),
			),
	)
	return err
}

func (m *Applet) DeleteByMd5(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("Md5").Eq(m.Md5),
					),
				),
				builder.Comment("Applet.DeleteByMd5"),
			),
	)
	return err
}

func (m *Applet) DeleteByAppletID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("AppletID").Eq(m.AppletID),
					),
				),
				builder.Comment("Applet.DeleteByAppletID"),
			),
	)
	return err
}
