// This is a generated source file. DO NOT EDIT
// Source: enums/cache_mode__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidCacheMode = errors.New("invalid CacheMode type")

func ParseCacheModeFromString(s string) (CacheMode, error) {
	switch s {
	default:
		return CACHE_MODE_UNKNOWN, InvalidCacheMode
	case "":
		return CACHE_MODE_UNKNOWN, nil
	case "MEMORY":
		return CACHE_MODE__MEMORY, nil
	case "REDIS":
		return CACHE_MODE__REDIS, nil
	}
}

func ParseCacheModeFromLabel(s string) (CacheMode, error) {
	switch s {
	default:
		return CACHE_MODE_UNKNOWN, InvalidCacheMode
	case "":
		return CACHE_MODE_UNKNOWN, nil
	case "MEMORY":
		return CACHE_MODE__MEMORY, nil
	case "REDIS":
		return CACHE_MODE__REDIS, nil
	}
}

func (v CacheMode) Int() int {
	return int(v)
}

func (v CacheMode) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case CACHE_MODE_UNKNOWN:
		return ""
	case CACHE_MODE__MEMORY:
		return "MEMORY"
	case CACHE_MODE__REDIS:
		return "REDIS"
	}
}

func (v CacheMode) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case CACHE_MODE_UNKNOWN:
		return ""
	case CACHE_MODE__MEMORY:
		return "MEMORY"
	case CACHE_MODE__REDIS:
		return "REDIS"
	}
}

func (v CacheMode) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.CacheMode"
}

func (v CacheMode) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{CACHE_MODE__MEMORY, CACHE_MODE__REDIS}
}

func (v CacheMode) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidCacheMode
	}
	return []byte(s), nil
}

func (v *CacheMode) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseCacheModeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *CacheMode) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = CacheMode(i)
	return nil
}

func (v CacheMode) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
