package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

// Device database model for device mangement
// @def primary                      ID
// @def unique_index UI_device_id    DeviceID
// @def unique_index UI_device_owner ProjectID DeviceMN
//
//go:generate toolkit gen model Device --database GwDB
type Device struct {
	datatypes.PrimaryID
	RelDevice
	DeviceInfo
}

type RelDevice struct {
	DeviceID types.SFID `db:"f_device_id" json:"deviceID"`
}

type DeviceInfo struct {
	ProjectID    types.SFID `db:"f_project_id"              json:"projectID"`
	DeviceMN     string     `db:"f_device_mn"               json:"deviceMN"`
	Manufacturer string     `db:"f_manufacturer,default=''" json:"manufacturer"`
}
