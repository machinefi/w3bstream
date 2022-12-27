// This is a generated source file. DO NOT EDIT
// Source: enums/config_type__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidConfigType = errors.New("invalid ConfigType type")

func ParseConfigTypeFromString(s string) (ConfigType, error) {
	switch s {
	default:
		return CONFIG_TYPE_UNKNOWN, InvalidConfigType
	case "":
		return CONFIG_TYPE_UNKNOWN, nil
	case "PROJECT_SCHEMA":
		return CONFIG_TYPE__PROJECT_SCHEMA, nil
	case "INSTANCE_CACHE":
		return CONFIG_TYPE__INSTANCE_CACHE, nil
	case "PROJECT_ENV":
		return CONFIG_TYPE__PROJECT_ENV, nil
	}
}

func ParseConfigTypeFromLabel(s string) (ConfigType, error) {
	switch s {
	default:
		return CONFIG_TYPE_UNKNOWN, InvalidConfigType
	case "":
		return CONFIG_TYPE_UNKNOWN, nil
	case "PROJECT_SCHEMA":
		return CONFIG_TYPE__PROJECT_SCHEMA, nil
	case "INSTANCE_CACHE":
		return CONFIG_TYPE__INSTANCE_CACHE, nil
	case "PROJECT_ENV":
		return CONFIG_TYPE__PROJECT_ENV, nil
	}
}

func (v ConfigType) Int() int {
	return int(v)
}

func (v ConfigType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case CONFIG_TYPE_UNKNOWN:
		return ""
	case CONFIG_TYPE__PROJECT_SCHEMA:
		return "PROJECT_SCHEMA"
	case CONFIG_TYPE__INSTANCE_CACHE:
		return "INSTANCE_CACHE"
	case CONFIG_TYPE__PROJECT_ENV:
		return "PROJECT_ENV"
	}
}

func (v ConfigType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case CONFIG_TYPE_UNKNOWN:
		return ""
	case CONFIG_TYPE__PROJECT_SCHEMA:
		return "PROJECT_SCHEMA"
	case CONFIG_TYPE__INSTANCE_CACHE:
		return "INSTANCE_CACHE"
	case CONFIG_TYPE__PROJECT_ENV:
		return "PROJECT_ENV"
	}
}

func (v ConfigType) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.ConfigType"
}

func (v ConfigType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{CONFIG_TYPE__PROJECT_SCHEMA, CONFIG_TYPE__INSTANCE_CACHE, CONFIG_TYPE__PROJECT_ENV}
}

func (v ConfigType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidConfigType
	}
	return []byte(s), nil
}

func (v *ConfigType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseConfigTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *ConfigType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = ConfigType(i)
	return nil
}

func (v ConfigType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
