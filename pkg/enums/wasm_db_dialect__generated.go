// This is a generated source file. DO NOT EDIT
// Source: enums/wasm_db_dialect__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidWasmDBDialect = errors.New("invalid WasmDBDialect type")

func ParseWasmDBDialectFromString(s string) (WasmDBDialect, error) {
	switch s {
	default:
		return WASM_DB_DIALECT_UNKNOWN, InvalidWasmDBDialect
	case "":
		return WASM_DB_DIALECT_UNKNOWN, nil
	case "POSTGRES":
		return WASM_DB_DIALECT__POSTGRES, nil
	}
}

func ParseWasmDBDialectFromLabel(s string) (WasmDBDialect, error) {
	switch s {
	default:
		return WASM_DB_DIALECT_UNKNOWN, InvalidWasmDBDialect
	case "":
		return WASM_DB_DIALECT_UNKNOWN, nil
	case "POSTGRES":
		return WASM_DB_DIALECT__POSTGRES, nil
	}
}

func (v WasmDBDialect) Int() int {
	return int(v)
}

func (v WasmDBDialect) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case WASM_DB_DIALECT_UNKNOWN:
		return ""
	case WASM_DB_DIALECT__POSTGRES:
		return "POSTGRES"
	}
}

func (v WasmDBDialect) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case WASM_DB_DIALECT_UNKNOWN:
		return ""
	case WASM_DB_DIALECT__POSTGRES:
		return "POSTGRES"
	}
}

func (v WasmDBDialect) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.WasmDBDialect"
}

func (v WasmDBDialect) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{WASM_DB_DIALECT__POSTGRES}
}

func (v WasmDBDialect) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidWasmDBDialect
	}
	return []byte(s), nil
}

func (v *WasmDBDialect) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseWasmDBDialectFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *WasmDBDialect) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = WasmDBDialect(i)
	return nil
}

func (v WasmDBDialect) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
