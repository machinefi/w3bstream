// This is a generated source file. DO NOT EDIT
// Source: schema/def_type__generated.go

package schema

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidDefType = errors.New("invalid DefType type")

func ParseDefTypeFromString(s string) (DefType, error) {
	switch s {
	default:
		return DEF_TYPE_UNKNOWN, InvalidDefType
	case "":
		return DEF_TYPE_UNKNOWN, nil
	case "PRIMARY":
		return DEF_TYPE__PRIMARY, nil
	case "INDEX":
		return DEF_TYPE__INDEX, nil
	case "UNIQUE_INDEX":
		return DEF_TYPE__UNIQUE_INDEX, nil
	}
}

func ParseDefTypeFromLabel(s string) (DefType, error) {
	switch s {
	default:
		return DEF_TYPE_UNKNOWN, InvalidDefType
	case "":
		return DEF_TYPE_UNKNOWN, nil
	case "PRIMARY":
		return DEF_TYPE__PRIMARY, nil
	case "INDEX":
		return DEF_TYPE__INDEX, nil
	case "UNIQUE_INDEX":
		return DEF_TYPE__UNIQUE_INDEX, nil
	}
}

func (v DefType) Int() int {
	return int(v)
}

func (v DefType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case DEF_TYPE_UNKNOWN:
		return ""
	case DEF_TYPE__PRIMARY:
		return "PRIMARY"
	case DEF_TYPE__INDEX:
		return "INDEX"
	case DEF_TYPE__UNIQUE_INDEX:
		return "UNIQUE_INDEX"
	}
}

func (v DefType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case DEF_TYPE_UNKNOWN:
		return ""
	case DEF_TYPE__PRIMARY:
		return "PRIMARY"
	case DEF_TYPE__INDEX:
		return "INDEX"
	case DEF_TYPE__UNIQUE_INDEX:
		return "UNIQUE_INDEX"
	}
}

func (v DefType) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/depends/schema.DefType"
}

func (v DefType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{DEF_TYPE__PRIMARY, DEF_TYPE__INDEX, DEF_TYPE__UNIQUE_INDEX}
}

func (v DefType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidDefType
	}
	return []byte(s), nil
}

func (v *DefType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseDefTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *DefType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = DefType(i)
	return nil
}

func (v DefType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
