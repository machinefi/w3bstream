// This is a generated source file. DO NOT EDIT
// Source: models/wasm_log__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var WasmLogTable *builder.Table

func init() {
	WasmLogTable = DB.Register(&WasmLog{})
}

type WasmLogIterator struct {
}

func (*WasmLogIterator) New() interface{} {
	return &WasmLog{}
}

func (*WasmLogIterator) Resolve(v interface{}) *WasmLog {
	return v.(*WasmLog)
}

func (*WasmLog) TableName() string {
	return "t_wasm_log"
}

func (*WasmLog) TableDesc() []string {
	return []string{
		"WasmLog database model event",
	}
}

func (*WasmLog) Comments() map[string]string {
	return map[string]string{}
}

func (*WasmLog) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*WasmLog) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*WasmLog) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *WasmLog) IndexFieldNames() []string {
	return []string{
		"ID",
		"WasmLogID",
	}
}

func (*WasmLog) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_wasm_log_id": []string{
			"WasmLogID",
		},
	}
}

func (*WasmLog) UniqueIndexUIWasmLogID() string {
	return "ui_wasm_log_id"
}

func (m *WasmLog) ColID() *builder.Column {
	return WasmLogTable.ColByFieldName(m.FieldID())
}

func (*WasmLog) FieldID() string {
	return "ID"
}

func (m *WasmLog) ColWasmLogID() *builder.Column {
	return WasmLogTable.ColByFieldName(m.FieldWasmLogID())
}

func (*WasmLog) FieldWasmLogID() string {
	return "WasmLogID"
}

func (m *WasmLog) ColProjectName() *builder.Column {
	return WasmLogTable.ColByFieldName(m.FieldProjectName())
}

func (*WasmLog) FieldProjectName() string {
	return "ProjectName"
}

func (m *WasmLog) ColAppletName() *builder.Column {
	return WasmLogTable.ColByFieldName(m.FieldAppletName())
}

func (*WasmLog) FieldAppletName() string {
	return "AppletName"
}

func (m *WasmLog) ColInstanceID() *builder.Column {
	return WasmLogTable.ColByFieldName(m.FieldInstanceID())
}

func (*WasmLog) FieldInstanceID() string {
	return "InstanceID"
}

func (m *WasmLog) ColLevel() *builder.Column {
	return WasmLogTable.ColByFieldName(m.FieldLevel())
}

func (*WasmLog) FieldLevel() string {
	return "Level"
}

func (m *WasmLog) ColLogTime() *builder.Column {
	return WasmLogTable.ColByFieldName(m.FieldLogTime())
}

func (*WasmLog) FieldLogTime() string {
	return "LogTime"
}

func (m *WasmLog) ColMsg() *builder.Column {
	return WasmLogTable.ColByFieldName(m.FieldMsg())
}

func (*WasmLog) FieldMsg() string {
	return "Msg"
}

func (m *WasmLog) ColCreatedAt() *builder.Column {
	return WasmLogTable.ColByFieldName(m.FieldCreatedAt())
}

func (*WasmLog) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *WasmLog) ColUpdatedAt() *builder.Column {
	return WasmLogTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*WasmLog) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *WasmLog) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *WasmLog) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *WasmLog) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]WasmLog, error) {
	var (
		tbl = db.T(m)
		lst = make([]WasmLog, 0)
	)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("WasmLog.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *WasmLog) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("WasmLog.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *WasmLog) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("WasmLog.FetchByID"),
			),
		m,
	)
	return err
}

func (m *WasmLog) FetchByWasmLogID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("WasmLogID").Eq(m.WasmLogID),
					),
				),
				builder.Comment("WasmLog.FetchByWasmLogID"),
			),
		m,
	)
	return err
}

func (m *WasmLog) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("WasmLog.UpdateByIDWithFVs"),
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

func (m *WasmLog) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *WasmLog) UpdateByWasmLogIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("WasmLogID").Eq(m.WasmLogID),
				),
				builder.Comment("WasmLog.UpdateByWasmLogIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByWasmLogID(db)
	}
	return nil
}

func (m *WasmLog) UpdateByWasmLogID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByWasmLogIDWithFVs(db, fvs)
}

func (m *WasmLog) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("WasmLog.Delete"),
			),
	)
	return err
}

func (m *WasmLog) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("WasmLog.DeleteByID"),
			),
	)
	return err
}

func (m *WasmLog) DeleteByWasmLogID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("WasmLogID").Eq(m.WasmLogID),
					),
				),
				builder.Comment("WasmLog.DeleteByWasmLogID"),
			),
	)
	return err
}
