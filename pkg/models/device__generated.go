// This is a generated source file. DO NOT EDIT
// Source: models/device__generated.go

package models

import (
	"fmt"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var DeviceTable *builder.Table

func init() {
	DeviceTable = GwDB.Register(&Device{})
}

type DeviceIterator struct {
}

func (*DeviceIterator) New() interface{} {
	return &Device{}
}

func (*DeviceIterator) Resolve(v interface{}) *Device {
	return v.(*Device)
}

func (*Device) TableName() string {
	return "t_device"
}

func (*Device) TableDesc() []string {
	return []string{
		"Device database model for device mangement",
	}
}

func (*Device) Comments() map[string]string {
	return map[string]string{}
}

func (*Device) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*Device) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*Device) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Device) IndexFieldNames() []string {
	return []string{
		"DeviceID",
		"DeviceMN",
		"ID",
		"ProjectID",
	}
}

func (*Device) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_device_id": []string{
			"DeviceID",
		},
		"ui_device_owner": []string{
			"ProjectID",
			"DeviceMN",
		},
	}
}

func (*Device) UniqueIndexUIDeviceID() string {
	return "ui_device_id"
}

func (*Device) UniqueIndexUIDeviceOwner() string {
	return "ui_device_owner"
}

func (m *Device) ColID() *builder.Column {
	return DeviceTable.ColByFieldName(m.FieldID())
}

func (*Device) FieldID() string {
	return "ID"
}

func (m *Device) ColDeviceID() *builder.Column {
	return DeviceTable.ColByFieldName(m.FieldDeviceID())
}

func (*Device) FieldDeviceID() string {
	return "DeviceID"
}

func (m *Device) ColProjectID() *builder.Column {
	return DeviceTable.ColByFieldName(m.FieldProjectID())
}

func (*Device) FieldProjectID() string {
	return "ProjectID"
}

func (m *Device) ColDeviceMN() *builder.Column {
	return DeviceTable.ColByFieldName(m.FieldDeviceMN())
}

func (*Device) FieldDeviceMN() string {
	return "DeviceMN"
}

func (m *Device) ColManufacturer() *builder.Column {
	return DeviceTable.ColByFieldName(m.FieldManufacturer())
}

func (*Device) FieldManufacturer() string {
	return "Manufacturer"
}

func (m *Device) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *Device) Create(db sqlx.DBExecutor) error {

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Device) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Device, error) {
	var (
		tbl = db.T(m)
		lst = make([]Device, 0)
	)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Device.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Device) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Device.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Device) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Device.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Device) FetchByDeviceID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("DeviceID").Eq(m.DeviceID),
					),
				),
				builder.Comment("Device.FetchByDeviceID"),
			),
		m,
	)
	return err
}

func (m *Device) FetchByProjectIDAndDeviceMN(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
						tbl.ColByFieldName("DeviceMN").Eq(m.DeviceMN),
					),
				),
				builder.Comment("Device.FetchByProjectIDAndDeviceMN"),
			),
		m,
	)
	return err
}

func (m *Device) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ID").Eq(m.ID),
				),
				builder.Comment("Device.UpdateByIDWithFVs"),
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

func (m *Device) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Device) UpdateByDeviceIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("DeviceID").Eq(m.DeviceID),
				),
				builder.Comment("Device.UpdateByDeviceIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByDeviceID(db)
	}
	return nil
}

func (m *Device) UpdateByDeviceID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByDeviceIDWithFVs(db, fvs)
}

func (m *Device) UpdateByProjectIDAndDeviceMNWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
					tbl.ColByFieldName("DeviceMN").Eq(m.DeviceMN),
				),
				builder.Comment("Device.UpdateByProjectIDAndDeviceMNWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByProjectIDAndDeviceMN(db)
	}
	return nil
}

func (m *Device) UpdateByProjectIDAndDeviceMN(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByProjectIDAndDeviceMNWithFVs(db, fvs)
}

func (m *Device) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Device.Delete"),
			),
	)
	return err
}

func (m *Device) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Device.DeleteByID"),
			),
	)
	return err
}

func (m *Device) DeleteByDeviceID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("DeviceID").Eq(m.DeviceID),
					),
				),
				builder.Comment("Device.DeleteByDeviceID"),
			),
	)
	return err
}

func (m *Device) DeleteByProjectIDAndDeviceMN(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
						tbl.ColByFieldName("DeviceMN").Eq(m.DeviceMN),
					),
				),
				builder.Comment("Device.DeleteByProjectIDAndDeviceMN"),
			),
	)
	return err
}
