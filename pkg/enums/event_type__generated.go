// This is a generated source file. DO NOT EDIT
// Source: enums/event_type__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/iotexproject/Bumblebee/kit/enum"
)

var InvalidEventType = errors.New("invalid EventType type")

func ParseEventTypeFromString(s string) (EventType, error) {
	switch s {
	default:
		return EVENT_TYPE_UNKNOWN, InvalidEventType
	case "":
		return EVENT_TYPE_UNKNOWN, nil
	case "ANY":
		return EVENT_TYPE__ANY, nil
	}
}

func ParseEventTypeFromLabel(s string) (EventType, error) {
	switch s {
	default:
		return EVENT_TYPE_UNKNOWN, InvalidEventType
	case "":
		return EVENT_TYPE_UNKNOWN, nil
	case "any event type":
		return EVENT_TYPE__ANY, nil
	}
}

func (v EventType) Int() int {
	return int(v)
}

func (v EventType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case EVENT_TYPE_UNKNOWN:
		return ""
	case EVENT_TYPE__ANY:
		return "ANY"
	}
}

func (v EventType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case EVENT_TYPE_UNKNOWN:
		return ""
	case EVENT_TYPE__ANY:
		return "any event type"
	}
}

func (v EventType) TypeName() string {
	return "github.com/iotexproject/w3bstream/pkg/enums.EventType"
}

func (v EventType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{EVENT_TYPE__ANY}
}

func (v EventType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidEventType
	}
	return []byte(s), nil
}

func (v *EventType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseEventTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *EventType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = EventType(i)
	return nil
}

func (v EventType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
