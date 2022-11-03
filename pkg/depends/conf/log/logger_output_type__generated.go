// This is a generated source file. DO NOT EDIT
// Source: log/logger_output_type__generated.go

package log

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidLoggerOutputType = errors.New("invalid LoggerOutputType type")

func ParseLoggerOutputTypeFromString(s string) (LoggerOutputType, error) {
	switch s {
	default:
		return LOGGER_OUTPUT_TYPE_UNKNOWN, InvalidLoggerOutputType
	case "":
		return LOGGER_OUTPUT_TYPE_UNKNOWN, nil
	case "ALWAYS":
		return LOGGER_OUTPUT_TYPE__ALWAYS, nil
	case "ON_FAILURE":
		return LOGGER_OUTPUT_TYPE__ON_FAILURE, nil
	case "NEVER":
		return LOGGER_OUTPUT_TYPE__NEVER, nil
	}
}

func ParseLoggerOutputTypeFromLabel(s string) (LoggerOutputType, error) {
	switch s {
	default:
		return LOGGER_OUTPUT_TYPE_UNKNOWN, InvalidLoggerOutputType
	case "":
		return LOGGER_OUTPUT_TYPE_UNKNOWN, nil
	case "ALWAYS":
		return LOGGER_OUTPUT_TYPE__ALWAYS, nil
	case "ON_FAILURE":
		return LOGGER_OUTPUT_TYPE__ON_FAILURE, nil
	case "NEVER":
		return LOGGER_OUTPUT_TYPE__NEVER, nil
	}
}

func (v LoggerOutputType) Int() int {
	return int(v)
}

func (v LoggerOutputType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case LOGGER_OUTPUT_TYPE_UNKNOWN:
		return ""
	case LOGGER_OUTPUT_TYPE__ALWAYS:
		return "ALWAYS"
	case LOGGER_OUTPUT_TYPE__ON_FAILURE:
		return "ON_FAILURE"
	case LOGGER_OUTPUT_TYPE__NEVER:
		return "NEVER"
	}
}

func (v LoggerOutputType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case LOGGER_OUTPUT_TYPE_UNKNOWN:
		return ""
	case LOGGER_OUTPUT_TYPE__ALWAYS:
		return "ALWAYS"
	case LOGGER_OUTPUT_TYPE__ON_FAILURE:
		return "ON_FAILURE"
	case LOGGER_OUTPUT_TYPE__NEVER:
		return "NEVER"
	}
}

func (v LoggerOutputType) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/depends/conf/log.LoggerOutputType"
}

func (v LoggerOutputType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{LOGGER_OUTPUT_TYPE__ALWAYS, LOGGER_OUTPUT_TYPE__ON_FAILURE, LOGGER_OUTPUT_TYPE__NEVER}
}

func (v LoggerOutputType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidLoggerOutputType
	}
	return []byte(s), nil
}

func (v *LoggerOutputType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseLoggerOutputTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *LoggerOutputType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = LoggerOutputType(i)
	return nil
}

func (v LoggerOutputType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
