// This is a generated source file. DO NOT EDIT
// Source: logger/format_type__generated.go

package logger

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidFormatType = errors.New("invalid FormatType type")

func ParseFormatTypeFromString(s string) (FormatType, error) {
	switch s {
	default:
		return FORMAT_TYPE_UNKNOWN, InvalidFormatType
	case "":
		return FORMAT_TYPE_UNKNOWN, nil
	case "JSON":
		return FORMAT_TYPE__JSON, nil
	case "TEXT":
		return FORMAT_TYPE__TEXT, nil
	}
}

func ParseFormatTypeFromLabel(s string) (FormatType, error) {
	switch s {
	default:
		return FORMAT_TYPE_UNKNOWN, InvalidFormatType
	case "":
		return FORMAT_TYPE_UNKNOWN, nil
	case "JSON":
		return FORMAT_TYPE__JSON, nil
	case "TEXT":
		return FORMAT_TYPE__TEXT, nil
	}
}

func (v FormatType) Int() int {
	return int(v)
}

func (v FormatType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case FORMAT_TYPE_UNKNOWN:
		return ""
	case FORMAT_TYPE__JSON:
		return "JSON"
	case FORMAT_TYPE__TEXT:
		return "TEXT"
	}
}

func (v FormatType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case FORMAT_TYPE_UNKNOWN:
		return ""
	case FORMAT_TYPE__JSON:
		return "JSON"
	case FORMAT_TYPE__TEXT:
		return "TEXT"
	}
}

func (v FormatType) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/depends/conf/logger.FormatType"
}

func (v FormatType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{FORMAT_TYPE__JSON, FORMAT_TYPE__TEXT}
}

func (v FormatType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidFormatType
	}
	return []byte(s), nil
}

func (v *FormatType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseFormatTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *FormatType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = FormatType(i)
	return nil
}

func (v FormatType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
