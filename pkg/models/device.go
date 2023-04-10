package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

// Device database model for device mangement
// @def primary                      ID
// @def unique_index UI_device_id    DeviceID
// @def unique_index UI_device_owner ProjectID SerialNumber Manufacturer
//
//go:generate toolkit gen model Device --database GwDB
type Device struct {
	datatypes.PrimaryID
	RelDevice
	DeviceInfo
	datatypes.OperationTimesWithDeleted
}

type RelDevice struct {
	DeviceID types.SFID `db:"f_device_id" json:"deviceID"`
}

type DeviceInfo struct {
	ProjectID    types.SFID `db:"f_project_id"              json:"projectID"`
	SerialNumber string     `db:"f_serial_number"           json:"serialNumber"`
	Manufacturer string     `db:"f_manufacturer,default=''" json:"manufacturer"`
}
