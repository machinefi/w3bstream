// This is a generated source file. DO NOT EDIT
// Source: enums/traffic_limit_type__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidTrafficLimitType = errors.New("invalid TrafficLimitType type")

func ParseTrafficLimitTypeFromString(s string) (TrafficLimitType, error) {
	switch s {
	default:
		return TRAFFIC_LIMIT_TYPE_UNKNOWN, InvalidTrafficLimitType
	case "":
		return TRAFFIC_LIMIT_TYPE_UNKNOWN, nil
	case "EVENT":
		return TRAFFIC_LIMIT_TYPE__EVENT, nil
	case "BLOCKCHAIN":
		return TRAFFIC_LIMIT_TYPE__BLOCKCHAIN, nil
	}
}

func ParseTrafficLimitTypeFromLabel(s string) (TrafficLimitType, error) {
	switch s {
	default:
		return TRAFFIC_LIMIT_TYPE_UNKNOWN, InvalidTrafficLimitType
	case "":
		return TRAFFIC_LIMIT_TYPE_UNKNOWN, nil
	case "EVENT":
		return TRAFFIC_LIMIT_TYPE__EVENT, nil
	case "BLOCKCHAIN":
		return TRAFFIC_LIMIT_TYPE__BLOCKCHAIN, nil
	}
}

func (v TrafficLimitType) Int() int {
	return int(v)
}

func (v TrafficLimitType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case TRAFFIC_LIMIT_TYPE_UNKNOWN:
		return ""
	case TRAFFIC_LIMIT_TYPE__EVENT:
		return "EVENT"
	case TRAFFIC_LIMIT_TYPE__BLOCKCHAIN:
		return "BLOCKCHAIN"
	}
}

func (v TrafficLimitType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case TRAFFIC_LIMIT_TYPE_UNKNOWN:
		return ""
	case TRAFFIC_LIMIT_TYPE__EVENT:
		return "EVENT"
	case TRAFFIC_LIMIT_TYPE__BLOCKCHAIN:
		return "BLOCKCHAIN"
	}
}

func (v TrafficLimitType) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.TrafficLimitType"
}

func (v TrafficLimitType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{TRAFFIC_LIMIT_TYPE__EVENT, TRAFFIC_LIMIT_TYPE__BLOCKCHAIN}
}

func (v TrafficLimitType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidTrafficLimitType
	}
	return []byte(s), nil
}

func (v *TrafficLimitType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseTrafficLimitTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *TrafficLimitType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = TrafficLimitType(i)
	return nil
}

func (v TrafficLimitType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
