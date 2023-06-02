// This is a generated source file. DO NOT EDIT
// Source: enums/traffic_cycle__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidTrafficCycle = errors.New("invalid TrafficCycle type")

func ParseTrafficCycleFromString(s string) (TrafficCycle, error) {
	switch s {
	default:
		return TRAFFIC_CYCLE_UNKNOWN, InvalidTrafficCycle
	case "":
		return TRAFFIC_CYCLE_UNKNOWN, nil
	case "MINUTE":
		return TRAFFIC_CYCLE__MINUTE, nil
	case "HOUR":
		return TRAFFIC_CYCLE__HOUR, nil
	case "DAY":
		return TRAFFIC_CYCLE__DAY, nil
	case "WEEK":
		return TRAFFIC_CYCLE__WEEK, nil
	case "MONTH":
		return TRAFFIC_CYCLE__MONTH, nil
	case "YEAR":
		return TRAFFIC_CYCLE__YEAR, nil
	}
}

func ParseTrafficCycleFromLabel(s string) (TrafficCycle, error) {
	switch s {
	default:
		return TRAFFIC_CYCLE_UNKNOWN, InvalidTrafficCycle
	case "":
		return TRAFFIC_CYCLE_UNKNOWN, nil
	case "MINUTE":
		return TRAFFIC_CYCLE__MINUTE, nil
	case "HOUR":
		return TRAFFIC_CYCLE__HOUR, nil
	case "DAY":
		return TRAFFIC_CYCLE__DAY, nil
	case "WEEK":
		return TRAFFIC_CYCLE__WEEK, nil
	case "MONTH":
		return TRAFFIC_CYCLE__MONTH, nil
	case "YEAR":
		return TRAFFIC_CYCLE__YEAR, nil
	}
}

func (v TrafficCycle) Int() int {
	return int(v)
}

func (v TrafficCycle) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case TRAFFIC_CYCLE_UNKNOWN:
		return ""
	case TRAFFIC_CYCLE__MINUTE:
		return "MINUTE"
	case TRAFFIC_CYCLE__HOUR:
		return "HOUR"
	case TRAFFIC_CYCLE__DAY:
		return "DAY"
	case TRAFFIC_CYCLE__WEEK:
		return "WEEK"
	case TRAFFIC_CYCLE__MONTH:
		return "MONTH"
	case TRAFFIC_CYCLE__YEAR:
		return "YEAR"
	}
}

func (v TrafficCycle) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case TRAFFIC_CYCLE_UNKNOWN:
		return ""
	case TRAFFIC_CYCLE__MINUTE:
		return "MINUTE"
	case TRAFFIC_CYCLE__HOUR:
		return "HOUR"
	case TRAFFIC_CYCLE__DAY:
		return "DAY"
	case TRAFFIC_CYCLE__WEEK:
		return "WEEK"
	case TRAFFIC_CYCLE__MONTH:
		return "MONTH"
	case TRAFFIC_CYCLE__YEAR:
		return "YEAR"
	}
}

func (v TrafficCycle) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.TrafficCycle"
}

func (v TrafficCycle) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{TRAFFIC_CYCLE__MINUTE, TRAFFIC_CYCLE__HOUR, TRAFFIC_CYCLE__DAY, TRAFFIC_CYCLE__WEEK, TRAFFIC_CYCLE__MONTH, TRAFFIC_CYCLE__YEAR}
}

func (v TrafficCycle) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidTrafficCycle
	}
	return []byte(s), nil
}

func (v *TrafficCycle) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseTrafficCycleFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *TrafficCycle) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = TrafficCycle(i)
	return nil
}

func (v TrafficCycle) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
