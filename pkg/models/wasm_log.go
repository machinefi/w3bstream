package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

// WasmLog database model event
// @def primary                     ID
// @def unique_index UI_wasm_log_id WasmLogID
// @def index        I_project_name ProjectName
// @def index        I_log_time     LogTime
// @def index        I_created_at   CreatedAt
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
	ProjectName string     `db:"f_project_name"               json:"projectName"`
	AppletName  string     `db:"f_applet_name,default=''"     json:"appletName"`
	InstanceID  types.SFID `db:"f_instance_id,default='0'"    json:"instanceID"`
	Src         string     `db:"f_src,default=''"             json:"src"`
	Level       string     `db:"f_level,default=''"           json:"level"`
	LogTime     int64      `db:"f_log_time,default='0'"       json:"logTime"`
	Msg         string     `db:"f_msg,default='',size=1024"   json:"msg"`
}

func BatchCreateWasmLogs(db sqlx.DBExecutor, vs ...*WasmLog) error {
	vals := make([]any, 0, len(vs)*8)
	for _, v := range vs {
		vals = append(vals, v.WasmLogID, v.ProjectName, v.AppletName, v.InstanceID, v.Src, v.Level, v.LogTime, v.Msg)
	}
	if len(vals) == 0 {
		return nil
	}

	m := &WasmLog{}
	t := db.T(m)

	_, err := db.Exec(builder.Insert().Into(t).Values(
		builder.Cols(
			m.FieldWasmLogID(),
			m.FieldProjectName(),
			m.FieldAppletName(),
			m.FieldInstanceID(),
			m.FieldSrc(),
			m.FieldLevel(),
			m.FieldLogTime(),
			m.FieldMsg(),
		),
		vals...,
	))
	return err
}
