package models

import (
	"database/sql/driver"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

var (
	DB        = sqlx.NewDatabase("demo").WithSchema("applet_management")
	MonitorDB = sqlx.NewDatabase("demo").WithSchema("monitor")
	WasmDB    = sqlx.NewDatabase("wasm_storage")
)

type Meta map[string]string

func (Meta) DataType(driver string) string { return "text" }

func (m Meta) Value() (driver.Value, error) { return datatypes.JSONValue(m) }

func (m *Meta) Scan(src interface{}) error { return datatypes.JSONScan(src, m) }

type Text string

func (Text) Datatype(driver string) string { return "text" }
