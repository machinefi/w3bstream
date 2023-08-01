// This is a generated source file. DO NOT EDIT
// Source: logger/output_type__generated.go

package logger

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidOutputType = errors.New("invalid OutputType type")

func ParseOutputTypeFromString(s string) (OutputType, error) {
	switch s {
	default:
		return OUTPUT_TYPE_UNKNOWN, InvalidOutputType
	case "":
		return OUTPUT_TYPE_UNKNOWN, nil
	case "ALWAYS":
		return OUTPUT_TYPE__ALWAYS, nil
	case "ON_FAILURE":
		return OUTPUT_TYPE__ON_FAILURE, nil
	case "NEVER":
		return OUTPUT_TYPE__NEVER, nil
	}
}

func ParseOutputTypeFromLabel(s string) (OutputType, error) {
	switch s {
	default:
		return OUTPUT_TYPE_UNKNOWN, InvalidOutputType
	case "":
		return OUTPUT_TYPE_UNKNOWN, nil
	case "ALWAYS":
		return OUTPUT_TYPE__ALWAYS, nil
	case "ON_FAILURE":
		return OUTPUT_TYPE__ON_FAILURE, nil
	case "NEVER":
		return OUTPUT_TYPE__NEVER, nil
	}
}

func (v OutputType) Int() int {
	return int(v)
}

func (v OutputType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case OUTPUT_TYPE_UNKNOWN:
		return ""
	case OUTPUT_TYPE__ALWAYS:
		return "ALWAYS"
	case OUTPUT_TYPE__ON_FAILURE:
		return "ON_FAILURE"
	case OUTPUT_TYPE__NEVER:
		return "NEVER"
	}
}

func (v OutputType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case OUTPUT_TYPE_UNKNOWN:
		return ""
	case OUTPUT_TYPE__ALWAYS:
		return "ALWAYS"
	case OUTPUT_TYPE__ON_FAILURE:
		return "ON_FAILURE"
	case OUTPUT_TYPE__NEVER:
		return "NEVER"
	}
}

func (v OutputType) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/depends/conf/logger.OutputType"
}

func (v OutputType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{OUTPUT_TYPE__ALWAYS, OUTPUT_TYPE__ON_FAILURE, OUTPUT_TYPE__NEVER}
}

func (v OutputType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidOutputType
	}
	return []byte(s), nil
}

func (v *OutputType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseOutputTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *OutputType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = OutputType(i)
	return nil
}

func (v OutputType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
