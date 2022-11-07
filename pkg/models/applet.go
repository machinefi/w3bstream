package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

// Applet database model applet
// @def primary                          ID
// @def unique_index UI_applet_id        AppletID
// @def unique_index UI_project_name     ProjectID Name
//
//go:generate toolkit gen model Applet --database DB
type Applet struct {
	datatypes.PrimaryID
	RelProject
	RelApplet
	RelResource
	AppletInfo
	datatypes.OperationTimes
}

type RelApplet struct {
	AppletID types.SFID `db:"f_applet_id" json:"appletID"`
}

type AppletInfo struct {
	Name     string `db:"f_name" json:"name"`
	WasmName string `db:"f_wasm_name" json:"wasmName"`
}
