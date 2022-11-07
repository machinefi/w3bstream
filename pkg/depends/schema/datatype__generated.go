// This is a generated source file. DO NOT EDIT
// Source: schema/datatype__generated.go

package schema

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidDatatype = errors.New("invalid Datatype type")

func ParseDatatypeFromString(s string) (Datatype, error) {
	switch s {
	default:
		return DATATYPE_UNKNOWN, InvalidDatatype
	case "":
		return DATATYPE_UNKNOWN, nil
	case "INT":
		return DATATYPE__INT, nil
	case "INT8":
		return DATATYPE__INT8, nil
	case "INT16":
		return DATATYPE__INT16, nil
	case "INT32":
		return DATATYPE__INT32, nil
	case "INT64":
		return DATATYPE__INT64, nil
	case "UINT":
		return DATATYPE__UINT, nil
	case "UINT8":
		return DATATYPE__UINT8, nil
	case "UINT16":
		return DATATYPE__UINT16, nil
	case "UINT32":
		return DATATYPE__UINT32, nil
	case "UINT64":
		return DATATYPE__UINT64, nil
	case "STRING":
		return DATATYPE__STRING, nil
	case "BOOL":
		return DATATYPE__BOOL, nil
	case "TIMESTAMP":
		return DATATYPE__TIMESTAMP, nil
	}
}

func ParseDatatypeFromLabel(s string) (Datatype, error) {
	switch s {
	default:
		return DATATYPE_UNKNOWN, InvalidDatatype
	case "":
		return DATATYPE_UNKNOWN, nil
	case "INT":
		return DATATYPE__INT, nil
	case "INT8":
		return DATATYPE__INT8, nil
	case "INT16":
		return DATATYPE__INT16, nil
	case "INT32":
		return DATATYPE__INT32, nil
	case "INT64":
		return DATATYPE__INT64, nil
	case "UINT":
		return DATATYPE__UINT, nil
	case "UINT8":
		return DATATYPE__UINT8, nil
	case "UINT16":
		return DATATYPE__UINT16, nil
	case "UINT32":
		return DATATYPE__UINT32, nil
	case "UINT64":
		return DATATYPE__UINT64, nil
	case "STRING":
		return DATATYPE__STRING, nil
	case "BOOL":
		return DATATYPE__BOOL, nil
	case "TIMESTAMP":
		return DATATYPE__TIMESTAMP, nil
	}
}

func (v Datatype) Int() int {
	return int(v)
}

func (v Datatype) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case DATATYPE_UNKNOWN:
		return ""
	case DATATYPE__INT:
		return "INT"
	case DATATYPE__INT8:
		return "INT8"
	case DATATYPE__INT16:
		return "INT16"
	case DATATYPE__INT32:
		return "INT32"
	case DATATYPE__INT64:
		return "INT64"
	case DATATYPE__UINT:
		return "UINT"
	case DATATYPE__UINT8:
		return "UINT8"
	case DATATYPE__UINT16:
		return "UINT16"
	case DATATYPE__UINT32:
		return "UINT32"
	case DATATYPE__UINT64:
		return "UINT64"
	case DATATYPE__STRING:
		return "STRING"
	case DATATYPE__BOOL:
		return "BOOL"
	case DATATYPE__TIMESTAMP:
		return "TIMESTAMP"
	}
}

func (v Datatype) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case DATATYPE_UNKNOWN:
		return ""
	case DATATYPE__INT:
		return "INT"
	case DATATYPE__INT8:
		return "INT8"
	case DATATYPE__INT16:
		return "INT16"
	case DATATYPE__INT32:
		return "INT32"
	case DATATYPE__INT64:
		return "INT64"
	case DATATYPE__UINT:
		return "UINT"
	case DATATYPE__UINT8:
		return "UINT8"
	case DATATYPE__UINT16:
		return "UINT16"
	case DATATYPE__UINT32:
		return "UINT32"
	case DATATYPE__UINT64:
		return "UINT64"
	case DATATYPE__STRING:
		return "STRING"
	case DATATYPE__BOOL:
		return "BOOL"
	case DATATYPE__TIMESTAMP:
		return "TIMESTAMP"
	}
}

func (v Datatype) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/depends/schema.Datatype"
}

func (v Datatype) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{DATATYPE__INT, DATATYPE__INT8, DATATYPE__INT16, DATATYPE__INT32, DATATYPE__INT64, DATATYPE__UINT, DATATYPE__UINT8, DATATYPE__UINT16, DATATYPE__UINT32, DATATYPE__UINT64, DATATYPE__STRING, DATATYPE__BOOL, DATATYPE__TIMESTAMP}
}

func (v Datatype) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidDatatype
	}
	return []byte(s), nil
}

func (v *Datatype) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseDatatypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *Datatype) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = Datatype(i)
	return nil
}

func (v Datatype) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
