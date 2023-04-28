package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

// WasmLog database model event
// @def primary                           ID
// @def unique_index UI_wasm_log_id       WasmLogID
//
//go:generate toolkit gen model WasmLog --database DB
type WasmLog struct {
	datatypes.PrimaryID
	RelWasmLog
	WasmLogInfo
	datatypes.OperationTimes
}

type RelWasmLog struct {
	WasmLogID types.SFID `db:"f_wasm_log_id" json:"wasmLogID"`
}

type WasmLogInfo struct {
	ProjectName string     `db:"f_project_name" json:"projectName"`
	AppletName  string     `db:"f_applet_name,default=''" json:"appletName"`
	InstanceID  types.SFID `db:"f_instance_id,default='0'" json:"instanceID"`
	Level       string     `db:"f_level,default=''" json:"level"`
	LogTime     int64      `db:"f_log_time,default='0'" json:"logTime"`
	Msg         string     `db:"f_msg,default='',size=1024" json:"msg"`
}
