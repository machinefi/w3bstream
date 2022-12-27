// This is a generated source file. DO NOT EDIT
// Source: models/config__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var ConfigTable *builder.Table

func init() {
	ConfigTable = DB.Register(&Config{})
}

type ConfigIterator struct {
}

func (*ConfigIterator) New() interface{} {
	return &Config{}
}

func (*ConfigIterator) Resolve(v interface{}) *Config {
	return v.(*Config)
}

func (*Config) TableName() string {
	return "t_config"
}

func (*Config) TableDesc() []string {
	return []string{
		"Config database model config for configuration management",
	}
}

func (*Config) Comments() map[string]string {
	return map[string]string{}
}

func (*Config) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*Config) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*Config) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Config) IndexFieldNames() []string {
	return []string{
		"ConfigID",
		"ID",
		"RelID",
		"Type",
	}
}

func (*Config) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_config_id": []string{
			"ConfigID",
		},
		"ui_rel_type": []string{
			"RelID",
			"Type",
		},
	}
}

func (*Config) UniqueIndexUIConfigID() string {
	return "ui_config_id"
}

func (*Config) UniqueIndexUIRelType() string {
	return "ui_rel_type"
}

func (m *Config) ColID() *builder.Column {
	return ConfigTable.ColByFieldName(m.FieldID())
}

func (*Config) FieldID() string {
	return "ID"
}

func (m *Config) ColConfigID() *builder.Column {
	return ConfigTable.ColByFieldName(m.FieldConfigID())
}

func (*Config) FieldConfigID() string {
	return "ConfigID"
}

func (m *Config) ColRelID() *builder.Column {
	return ConfigTable.ColByFieldName(m.FieldRelID())
}

func (*Config) FieldRelID() string {
	return "RelID"
}

func (m *Config) ColType() *builder.Column {
	return ConfigTable.ColByFieldName(m.FieldType())
}

func (*Config) FieldType() string {
	return "Type"
}

func (m *Config) ColValue() *builder.Column {
	return ConfigTable.ColByFieldName(m.FieldValue())
}

func (*Config) FieldValue() string {
	return "Value"
}

func (m *Config) ColCreatedAt() *builder.Column {
	return ConfigTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Config) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Config) ColUpdatedAt() *builder.Column {
	return ConfigTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Config) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Config) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *Config) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Config) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Config, error) {
	var (
		tbl = db.T(m)
		lst = make([]Config, 0)
	)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Config.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Config) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Config.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Config) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Config.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Config) FetchByConfigID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ConfigID").Eq(m.ConfigID),
					),
				),
				builder.Comment("Config.FetchByConfigID"),
			),
		m,
	)
	return err
}

func (m *Config) FetchByRelIDAndType(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("RelID").Eq(m.RelID),
						tbl.ColByFieldName("Type").Eq(m.Type),
					),
				),
				builder.Comment("Config.FetchByRelIDAndType"),
			),
		m,
	)
	return err
}

func (m *Config) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("Config.UpdateByIDWithFVs"),
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

func (m *Config) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Config) UpdateByConfigIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ConfigID").Eq(m.ConfigID),
				),
				builder.Comment("Config.UpdateByConfigIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByConfigID(db)
	}
	return nil
}

func (m *Config) UpdateByConfigID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByConfigIDWithFVs(db, fvs)
}

func (m *Config) UpdateByRelIDAndTypeWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("RelID").Eq(m.RelID),
					tbl.ColByFieldName("Type").Eq(m.Type),
				),
				builder.Comment("Config.UpdateByRelIDAndTypeWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByRelIDAndType(db)
	}
	return nil
}

func (m *Config) UpdateByRelIDAndType(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByRelIDAndTypeWithFVs(db, fvs)
}

func (m *Config) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Config.Delete"),
			),
	)
	return err
}

func (m *Config) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Config.DeleteByID"),
			),
	)
	return err
}

func (m *Config) DeleteByConfigID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ConfigID").Eq(m.ConfigID),
					),
				),
				builder.Comment("Config.DeleteByConfigID"),
			),
	)
	return err
}

func (m *Config) DeleteByRelIDAndType(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("RelID").Eq(m.RelID),
						tbl.ColByFieldName("Type").Eq(m.Type),
					),
				),
				builder.Comment("Config.DeleteByRelIDAndType"),
			),
	)
	return err
}
