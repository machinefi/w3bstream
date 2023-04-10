// This is a generated source file. DO NOT EDIT
// Source: enums/proto_mode__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidProtoMode = errors.New("invalid ProtoMode type")

func ParseProtoModeFromString(s string) (ProtoMode, error) {
	switch s {
	default:
		return PROTO_MODE_UNKNOWN, InvalidProtoMode
	case "":
		return PROTO_MODE_UNKNOWN, nil
	case "PEBBLE":
		return PROTO_MODE__PEBBLE, nil
	}
}

func ParseProtoModeFromLabel(s string) (ProtoMode, error) {
	switch s {
	default:
		return PROTO_MODE_UNKNOWN, InvalidProtoMode
	case "":
		return PROTO_MODE_UNKNOWN, nil
	case "PEBBLE":
		return PROTO_MODE__PEBBLE, nil
	}
}

func (v ProtoMode) Int() int {
	return int(v)
}

func (v ProtoMode) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case PROTO_MODE_UNKNOWN:
		return ""
	case PROTO_MODE__PEBBLE:
		return "PEBBLE"
	}
}

func (v ProtoMode) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case PROTO_MODE_UNKNOWN:
		return ""
	case PROTO_MODE__PEBBLE:
		return "PEBBLE"
	}
}

func (v ProtoMode) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.ProtoMode"
}

func (v ProtoMode) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{PROTO_MODE__PEBBLE}
}

func (v ProtoMode) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidProtoMode
	}
	return []byte(s), nil
}

func (v *ProtoMode) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseProtoModeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *ProtoMode) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = ProtoMode(i)
	return nil
}

func (v ProtoMode) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
