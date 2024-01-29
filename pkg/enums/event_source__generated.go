// This is a generated source file. DO NOT EDIT
// Source: enums/event_source__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidEventSource = errors.New("invalid EventSource type")

func ParseEventSourceFromString(s string) (EventSource, error) {
	switch s {
	default:
		return EVENT_SOURCE_UNKNOWN, InvalidEventSource
	case "":
		return EVENT_SOURCE_UNKNOWN, nil
	case "MQTT":
		return EVENT_SOURCE__MQTT, nil
	case "HTTP":
		return EVENT_SOURCE__HTTP, nil
	}
}

func ParseEventSourceFromLabel(s string) (EventSource, error) {
	switch s {
	default:
		return EVENT_SOURCE_UNKNOWN, InvalidEventSource
	case "":
		return EVENT_SOURCE_UNKNOWN, nil
	case "MQTT":
		return EVENT_SOURCE__MQTT, nil
	case "HTTP":
		return EVENT_SOURCE__HTTP, nil
	}
}

func (v EventSource) Int() int {
	return int(v)
}

func (v EventSource) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case EVENT_SOURCE_UNKNOWN:
		return ""
	case EVENT_SOURCE__MQTT:
		return "MQTT"
	case EVENT_SOURCE__HTTP:
		return "HTTP"
	}
}

func (v EventSource) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case EVENT_SOURCE_UNKNOWN:
		return ""
	case EVENT_SOURCE__MQTT:
		return "MQTT"
	case EVENT_SOURCE__HTTP:
		return "HTTP"
	}
}

func (v EventSource) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.EventSource"
}

func (v EventSource) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{EVENT_SOURCE__MQTT, EVENT_SOURCE__HTTP}
}

func (v EventSource) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidEventSource
	}
	return []byte(s), nil
}

func (v *EventSource) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseEventSourceFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *EventSource) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = EventSource(i)
	return nil
}

func (v EventSource) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
