// This is a generated source file. DO NOT EDIT
// Source: wasm/log_level__generated.go

package consts

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidLogLevel = errors.New("invalid LogLevel type")

func ParseLogLevelFromString(s string) (LogLevel, error) {
	switch s {
	default:
		return LOG_LEVEL_UNKNOWN, InvalidLogLevel
	case "":
		return LOG_LEVEL_UNKNOWN, nil
	case "ERROR":
		return LOG_LEVEL__ERROR, nil
	case "WARN":
		return LOG_LEVEL__WARN, nil
	case "INFO":
		return LOG_LEVEL__INFO, nil
	case "DEBUG":
		return LOG_LEVEL__DEBUG, nil
	}
}

func ParseLogLevelFromLabel(s string) (LogLevel, error) {
	switch s {
	default:
		return LOG_LEVEL_UNKNOWN, InvalidLogLevel
	case "":
		return LOG_LEVEL_UNKNOWN, nil
	case "ERROR":
		return LOG_LEVEL__ERROR, nil
	case "WARN":
		return LOG_LEVEL__WARN, nil
	case "INFO":
		return LOG_LEVEL__INFO, nil
	case "DEBUG":
		return LOG_LEVEL__DEBUG, nil
	}
}

func (v LogLevel) Int() int {
	return int(v)
}

func (v LogLevel) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case LOG_LEVEL_UNKNOWN:
		return ""
	case LOG_LEVEL__ERROR:
		return "ERROR"
	case LOG_LEVEL__WARN:
		return "WARN"
	case LOG_LEVEL__INFO:
		return "INFO"
	case LOG_LEVEL__DEBUG:
		return "DEBUG"
	}
}

func (v LogLevel) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case LOG_LEVEL_UNKNOWN:
		return ""
	case LOG_LEVEL__ERROR:
		return "ERROR"
	case LOG_LEVEL__WARN:
		return "WARN"
	case LOG_LEVEL__INFO:
		return "INFO"
	case LOG_LEVEL__DEBUG:
		return "DEBUG"
	}
}

func (v LogLevel) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/modules/wasm.LogLevel"
}

func (v LogLevel) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{LOG_LEVEL__ERROR, LOG_LEVEL__WARN, LOG_LEVEL__INFO, LOG_LEVEL__DEBUG}
}

func (v LogLevel) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidLogLevel
	}
	return []byte(s), nil
}

func (v *LogLevel) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseLogLevelFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *LogLevel) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = LogLevel(i)
	return nil
}

func (v LogLevel) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
