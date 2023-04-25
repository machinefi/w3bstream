// This is a generated source file. DO NOT EDIT
// Source: enums/wasm_db_datatype__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidWasmDBDatatype = errors.New("invalid WasmDBDatatype type")

func ParseWasmDBDatatypeFromString(s string) (WasmDBDatatype, error) {
	switch s {
	default:
		return WASM_DB_DATATYPE_UNKNOWN, InvalidWasmDBDatatype
	case "":
		return WASM_DB_DATATYPE_UNKNOWN, nil
	case "INT":
		return WASM_DB_DATATYPE__INT, nil
	case "INT8":
		return WASM_DB_DATATYPE__INT8, nil
	case "INT16":
		return WASM_DB_DATATYPE__INT16, nil
	case "INT32":
		return WASM_DB_DATATYPE__INT32, nil
	case "INT64":
		return WASM_DB_DATATYPE__INT64, nil
	case "UINT":
		return WASM_DB_DATATYPE__UINT, nil
	case "UINT8":
		return WASM_DB_DATATYPE__UINT8, nil
	case "UINT16":
		return WASM_DB_DATATYPE__UINT16, nil
	case "UINT32":
		return WASM_DB_DATATYPE__UINT32, nil
	case "UINT64":
		return WASM_DB_DATATYPE__UINT64, nil
	case "FLOAT32":
		return WASM_DB_DATATYPE__FLOAT32, nil
	case "FLOAT64":
		return WASM_DB_DATATYPE__FLOAT64, nil
	case "TEXT":
		return WASM_DB_DATATYPE__TEXT, nil
	case "BOOL":
		return WASM_DB_DATATYPE__BOOL, nil
	case "TIMESTAMP":
		return WASM_DB_DATATYPE__TIMESTAMP, nil
	case "DECIMAL":
		return WASM_DB_DATATYPE__DECIMAL, nil
	case "NUMERIC":
		return WASM_DB_DATATYPE__NUMERIC, nil
	}
}

func ParseWasmDBDatatypeFromLabel(s string) (WasmDBDatatype, error) {
	switch s {
	default:
		return WASM_DB_DATATYPE_UNKNOWN, InvalidWasmDBDatatype
	case "":
		return WASM_DB_DATATYPE_UNKNOWN, nil
	case "INT":
		return WASM_DB_DATATYPE__INT, nil
	case "INT8":
		return WASM_DB_DATATYPE__INT8, nil
	case "INT16":
		return WASM_DB_DATATYPE__INT16, nil
	case "INT32":
		return WASM_DB_DATATYPE__INT32, nil
	case "INT64":
		return WASM_DB_DATATYPE__INT64, nil
	case "UINT":
		return WASM_DB_DATATYPE__UINT, nil
	case "UINT8":
		return WASM_DB_DATATYPE__UINT8, nil
	case "UINT16":
		return WASM_DB_DATATYPE__UINT16, nil
	case "UINT32":
		return WASM_DB_DATATYPE__UINT32, nil
	case "UINT64":
		return WASM_DB_DATATYPE__UINT64, nil
	case "FLOAT32":
		return WASM_DB_DATATYPE__FLOAT32, nil
	case "FLOAT64":
		return WASM_DB_DATATYPE__FLOAT64, nil
	case "TEXT":
		return WASM_DB_DATATYPE__TEXT, nil
	case "BOOL":
		return WASM_DB_DATATYPE__BOOL, nil
	case "use epoch timestamp (integer, UTC)":
		return WASM_DB_DATATYPE__TIMESTAMP, nil
	case "DECIMAL":
		return WASM_DB_DATATYPE__DECIMAL, nil
	case "NUMERIC":
		return WASM_DB_DATATYPE__NUMERIC, nil
	}
}

func (v WasmDBDatatype) Int() int {
	return int(v)
}

func (v WasmDBDatatype) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case WASM_DB_DATATYPE_UNKNOWN:
		return ""
	case WASM_DB_DATATYPE__INT:
		return "INT"
	case WASM_DB_DATATYPE__INT8:
		return "INT8"
	case WASM_DB_DATATYPE__INT16:
		return "INT16"
	case WASM_DB_DATATYPE__INT32:
		return "INT32"
	case WASM_DB_DATATYPE__INT64:
		return "INT64"
	case WASM_DB_DATATYPE__UINT:
		return "UINT"
	case WASM_DB_DATATYPE__UINT8:
		return "UINT8"
	case WASM_DB_DATATYPE__UINT16:
		return "UINT16"
	case WASM_DB_DATATYPE__UINT32:
		return "UINT32"
	case WASM_DB_DATATYPE__UINT64:
		return "UINT64"
	case WASM_DB_DATATYPE__FLOAT32:
		return "FLOAT32"
	case WASM_DB_DATATYPE__FLOAT64:
		return "FLOAT64"
	case WASM_DB_DATATYPE__TEXT:
		return "TEXT"
	case WASM_DB_DATATYPE__BOOL:
		return "BOOL"
	case WASM_DB_DATATYPE__TIMESTAMP:
		return "TIMESTAMP"
	case WASM_DB_DATATYPE__DECIMAL:
		return "DECIMAL"
	case WASM_DB_DATATYPE__NUMERIC:
		return "NUMERIC"
	}
}

func (v WasmDBDatatype) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case WASM_DB_DATATYPE_UNKNOWN:
		return ""
	case WASM_DB_DATATYPE__INT:
		return "INT"
	case WASM_DB_DATATYPE__INT8:
		return "INT8"
	case WASM_DB_DATATYPE__INT16:
		return "INT16"
	case WASM_DB_DATATYPE__INT32:
		return "INT32"
	case WASM_DB_DATATYPE__INT64:
		return "INT64"
	case WASM_DB_DATATYPE__UINT:
		return "UINT"
	case WASM_DB_DATATYPE__UINT8:
		return "UINT8"
	case WASM_DB_DATATYPE__UINT16:
		return "UINT16"
	case WASM_DB_DATATYPE__UINT32:
		return "UINT32"
	case WASM_DB_DATATYPE__UINT64:
		return "UINT64"
	case WASM_DB_DATATYPE__FLOAT32:
		return "FLOAT32"
	case WASM_DB_DATATYPE__FLOAT64:
		return "FLOAT64"
	case WASM_DB_DATATYPE__TEXT:
		return "TEXT"
	case WASM_DB_DATATYPE__BOOL:
		return "BOOL"
	case WASM_DB_DATATYPE__TIMESTAMP:
		return "use epoch timestamp (integer, UTC)"
	case WASM_DB_DATATYPE__DECIMAL:
		return "DECIMAL"
	case WASM_DB_DATATYPE__NUMERIC:
		return "NUMERIC"
	}
}

func (v WasmDBDatatype) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.WasmDBDatatype"
}

func (v WasmDBDatatype) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{WASM_DB_DATATYPE__INT, WASM_DB_DATATYPE__INT8, WASM_DB_DATATYPE__INT16, WASM_DB_DATATYPE__INT32, WASM_DB_DATATYPE__INT64, WASM_DB_DATATYPE__UINT, WASM_DB_DATATYPE__UINT8, WASM_DB_DATATYPE__UINT16, WASM_DB_DATATYPE__UINT32, WASM_DB_DATATYPE__UINT64, WASM_DB_DATATYPE__FLOAT32, WASM_DB_DATATYPE__FLOAT64, WASM_DB_DATATYPE__TEXT, WASM_DB_DATATYPE__BOOL, WASM_DB_DATATYPE__TIMESTAMP, WASM_DB_DATATYPE__DECIMAL, WASM_DB_DATATYPE__NUMERIC}
}

func (v WasmDBDatatype) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidWasmDBDatatype
	}
	return []byte(s), nil
}

func (v *WasmDBDatatype) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseWasmDBDatatypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *WasmDBDatatype) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = WasmDBDatatype(i)
	return nil
}

func (v WasmDBDatatype) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
