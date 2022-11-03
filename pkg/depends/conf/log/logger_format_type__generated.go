// This is a generated source file. DO NOT EDIT
// Source: log/logger_format_type__generated.go

package log

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidLoggerFormatType = errors.New("invalid LoggerFormatType type")

func ParseLoggerFormatTypeFromString(s string) (LoggerFormatType, error) {
	switch s {
	default:
		return LOGGER_FORMAT_TYPE_UNKNOWN, InvalidLoggerFormatType
	case "":
		return LOGGER_FORMAT_TYPE_UNKNOWN, nil
	case "JSON":
		return LOGGER_FORMAT_TYPE__JSON, nil
	case "TEXT":
		return LOGGER_FORMAT_TYPE__TEXT, nil
	}
}

func ParseLoggerFormatTypeFromLabel(s string) (LoggerFormatType, error) {
	switch s {
	default:
		return LOGGER_FORMAT_TYPE_UNKNOWN, InvalidLoggerFormatType
	case "":
		return LOGGER_FORMAT_TYPE_UNKNOWN, nil
	case "JSON":
		return LOGGER_FORMAT_TYPE__JSON, nil
	case "TEXT":
		return LOGGER_FORMAT_TYPE__TEXT, nil
	}
}

func (v LoggerFormatType) Int() int {
	return int(v)
}

func (v LoggerFormatType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case LOGGER_FORMAT_TYPE_UNKNOWN:
		return ""
	case LOGGER_FORMAT_TYPE__JSON:
		return "JSON"
	case LOGGER_FORMAT_TYPE__TEXT:
		return "TEXT"
	}
}

func (v LoggerFormatType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case LOGGER_FORMAT_TYPE_UNKNOWN:
		return ""
	case LOGGER_FORMAT_TYPE__JSON:
		return "JSON"
	case LOGGER_FORMAT_TYPE__TEXT:
		return "TEXT"
	}
}

func (v LoggerFormatType) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/depends/conf/log.LoggerFormatType"
}

func (v LoggerFormatType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{LOGGER_FORMAT_TYPE__JSON, LOGGER_FORMAT_TYPE__TEXT}
}

func (v LoggerFormatType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidLoggerFormatType
	}
	return []byte(s), nil
}

func (v *LoggerFormatType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseLoggerFormatTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *LoggerFormatType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = LoggerFormatType(i)
	return nil
}

func (v LoggerFormatType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
