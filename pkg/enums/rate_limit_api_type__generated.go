// This is a generated source file. DO NOT EDIT
// Source: enums/rate_limit_api_type__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidRateLimitApiType = errors.New("invalid RateLimitApiType type")

func ParseRateLimitApiTypeFromString(s string) (RateLimitApiType, error) {
	switch s {
	default:
		return RATE_LIMIT_API_TYPE_UNKNOWN, InvalidRateLimitApiType
	case "":
		return RATE_LIMIT_API_TYPE_UNKNOWN, nil
	case "EVENT":
		return RATE_LIMIT_API_TYPE__EVENT, nil
	case "BLOCKCHAIN":
		return RATE_LIMIT_API_TYPE__BLOCKCHAIN, nil
	}
}

func ParseRateLimitApiTypeFromLabel(s string) (RateLimitApiType, error) {
	switch s {
	default:
		return RATE_LIMIT_API_TYPE_UNKNOWN, InvalidRateLimitApiType
	case "":
		return RATE_LIMIT_API_TYPE_UNKNOWN, nil
	case "EVENT":
		return RATE_LIMIT_API_TYPE__EVENT, nil
	case "BLOCKCHAIN":
		return RATE_LIMIT_API_TYPE__BLOCKCHAIN, nil
	}
}

func (v RateLimitApiType) Int() int {
	return int(v)
}

func (v RateLimitApiType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case RATE_LIMIT_API_TYPE_UNKNOWN:
		return ""
	case RATE_LIMIT_API_TYPE__EVENT:
		return "EVENT"
	case RATE_LIMIT_API_TYPE__BLOCKCHAIN:
		return "BLOCKCHAIN"
	}
}

func (v RateLimitApiType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case RATE_LIMIT_API_TYPE_UNKNOWN:
		return ""
	case RATE_LIMIT_API_TYPE__EVENT:
		return "EVENT"
	case RATE_LIMIT_API_TYPE__BLOCKCHAIN:
		return "BLOCKCHAIN"
	}
}

func (v RateLimitApiType) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.RateLimitApiType"
}

func (v RateLimitApiType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{RATE_LIMIT_API_TYPE__EVENT, RATE_LIMIT_API_TYPE__BLOCKCHAIN}
}

func (v RateLimitApiType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidRateLimitApiType
	}
	return []byte(s), nil
}

func (v *RateLimitApiType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseRateLimitApiTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *RateLimitApiType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = RateLimitApiType(i)
	return nil
}

func (v RateLimitApiType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
