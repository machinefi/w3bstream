// This is a generated source file. DO NOT EDIT
// Source: models/instance__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/sqlx"
	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"
)

var InstanceTable *builder.Table

func init() {
	InstanceTable = DB.Register(&Instance{})
}

type InstanceIterator struct {
}

func (*InstanceIterator) New() interface{} {
	return &Instance{}
}

func (*InstanceIterator) Resolve(v interface{}) *Instance {
	return v.(*Instance)
}

func (*Instance) TableName() string {
	return "t_instance"
}

func (*Instance) TableDesc() []string {
	return []string{
		"Instance database model instance",
	}
}

func (*Instance) Comments() map[string]string {
	return map[string]string{}
}

func (*Instance) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*Instance) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*Instance) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (*Instance) Indexes() builder.Indexes {
	return builder.Indexes{
		"i_applet_id": []string{
			"AppletID",
		},
		"i_path": []string{
			"Path",
		},
	}
}

func (m *Instance) IndexFieldNames() []string {
	return []string{
		"AppletID",
		"ID",
		"InstanceID",
		"Path",
	}
}

func (*Instance) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_instance_id": []string{
			"InstanceID",
			"DeletedAt",
		},
	}
}

func (*Instance) UniqueIndexUIInstanceID() string {
	return "ui_instance_id"
}

func (m *Instance) ColID() *builder.Column {
	return InstanceTable.ColByFieldName(m.FieldID())
}

func (*Instance) FieldID() string {
	return "ID"
}

func (m *Instance) ColInstanceID() *builder.Column {
	return InstanceTable.ColByFieldName(m.FieldInstanceID())
}

func (*Instance) FieldInstanceID() string {
	return "InstanceID"
}

func (m *Instance) ColAppletID() *builder.Column {
	return InstanceTable.ColByFieldName(m.FieldAppletID())
}

func (*Instance) FieldAppletID() string {
	return "AppletID"
}

func (m *Instance) ColPath() *builder.Column {
	return InstanceTable.ColByFieldName(m.FieldPath())
}

func (*Instance) FieldPath() string {
	return "Path"
}

func (m *Instance) ColState() *builder.Column {
	return InstanceTable.ColByFieldName(m.FieldState())
}

func (*Instance) FieldState() string {
	return "State"
}

func (m *Instance) ColCreatedAt() *builder.Column {
	return InstanceTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Instance) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Instance) ColUpdatedAt() *builder.Column {
	return InstanceTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Instance) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Instance) ColDeletedAt() *builder.Column {
	return InstanceTable.ColByFieldName(m.FieldDeletedAt())
}

func (*Instance) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *Instance) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *Instance) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Instance) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Instance, error) {
	var (
		tbl = db.T(m)
		lst = make([]Instance, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Instance.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Instance) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Instance.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Instance) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Instance.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Instance) FetchByInstanceID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("InstanceID").Eq(m.InstanceID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Instance.FetchByInstanceID"),
			),
		m,
	)
	return err
}

func (m *Instance) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("Instance.UpdateByIDWithFVs"),
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

func (m *Instance) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Instance) UpdateByInstanceIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("InstanceID").Eq(m.InstanceID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Instance.UpdateByInstanceIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByInstanceID(db)
	}
	return nil
}

func (m *Instance) UpdateByInstanceID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByInstanceIDWithFVs(db, fvs)
}

func (m *Instance) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Instance.Delete"),
			),
	)
	return err
}

func (m *Instance) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Instance.DeleteByID"),
			),
	)
	return err
}

func (m *Instance) SoftDeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Instance.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *Instance) DeleteByInstanceID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("InstanceID").Eq(m.InstanceID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Instance.DeleteByInstanceID"),
			),
	)
	return err
}

func (m *Instance) SoftDeleteByInstanceID(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("InstanceID").Eq(m.InstanceID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Instance.SoftDeleteByInstanceID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
