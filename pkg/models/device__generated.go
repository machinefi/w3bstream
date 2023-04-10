// This is a generated source file. DO NOT EDIT
// Source: models/device__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
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
		"ID",
		"Manufacturer",
		"ProjectID",
		"SerialNumber",
	}
}

func (*Device) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_device_id": []string{
			"DeviceID",
			"DeletedAt",
		},
		"ui_device_owner": []string{
			"ProjectID",
			"SerialNumber",
			"Manufacturer",
			"DeletedAt",
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

func (m *Device) ColSerialNumber() *builder.Column {
	return DeviceTable.ColByFieldName(m.FieldSerialNumber())
}

func (*Device) FieldSerialNumber() string {
	return "SerialNumber"
}

func (m *Device) ColManufacturer() *builder.Column {
	return DeviceTable.ColByFieldName(m.FieldManufacturer())
}

func (*Device) FieldManufacturer() string {
	return "Manufacturer"
}

func (m *Device) ColCreatedAt() *builder.Column {
	return DeviceTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Device) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Device) ColUpdatedAt() *builder.Column {
	return DeviceTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Device) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Device) ColDeletedAt() *builder.Column {
	return DeviceTable.ColByFieldName(m.FieldDeletedAt())
}

func (*Device) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *Device) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *Device) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Device) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Device, error) {
	var (
		tbl = db.T(m)
		lst = make([]Device, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Device.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Device) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
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
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
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
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Device.FetchByDeviceID"),
			),
		m,
	)
	return err
}

func (m *Device) FetchByProjectIDAndSerialNumberAndManufacturer(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
						tbl.ColByFieldName("SerialNumber").Eq(m.SerialNumber),
						tbl.ColByFieldName("Manufacturer").Eq(m.Manufacturer),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Device.FetchByProjectIDAndSerialNumberAndManufacturer"),
			),
		m,
	)
	return err
}

func (m *Device) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("DeviceID").Eq(m.DeviceID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
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

func (m *Device) UpdateByProjectIDAndSerialNumberAndManufacturerWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
					tbl.ColByFieldName("SerialNumber").Eq(m.SerialNumber),
					tbl.ColByFieldName("Manufacturer").Eq(m.Manufacturer),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Device.UpdateByProjectIDAndSerialNumberAndManufacturerWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByProjectIDAndSerialNumberAndManufacturer(db)
	}
	return nil
}

func (m *Device) UpdateByProjectIDAndSerialNumberAndManufacturer(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByProjectIDAndSerialNumberAndManufacturerWithFVs(db, fvs)
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
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Device.DeleteByID"),
			),
	)
	return err
}

func (m *Device) SoftDeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Device.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
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
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Device.DeleteByDeviceID"),
			),
	)
	return err
}

func (m *Device) SoftDeleteByDeviceID(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("DeviceID").Eq(m.DeviceID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Device.SoftDeleteByDeviceID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *Device) DeleteByProjectIDAndSerialNumberAndManufacturer(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
						tbl.ColByFieldName("SerialNumber").Eq(m.SerialNumber),
						tbl.ColByFieldName("Manufacturer").Eq(m.Manufacturer),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Device.DeleteByProjectIDAndSerialNumberAndManufacturer"),
			),
	)
	return err
}

func (m *Device) SoftDeleteByProjectIDAndSerialNumberAndManufacturer(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("ProjectID").Eq(m.ProjectID),
					tbl.ColByFieldName("SerialNumber").Eq(m.SerialNumber),
					tbl.ColByFieldName("Manufacturer").Eq(m.Manufacturer),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Device.SoftDeleteByProjectIDAndSerialNumberAndManufacturer"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
