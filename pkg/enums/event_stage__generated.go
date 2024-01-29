// This is a generated source file. DO NOT EDIT
// Source: enums/event_stage__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidEventStage = errors.New("invalid EventStage type")

func ParseEventStageFromString(s string) (EventStage, error) {
	switch s {
	default:
		return EVENT_STAGE_UNKNOWN, InvalidEventStage
	case "":
		return EVENT_STAGE_UNKNOWN, nil
	case "RECEIVED":
		return EVENT_STAGE__RECEIVED, nil
	case "HANDLED":
		return EVENT_STAGE__HANDLED, nil
	case "COMPLETED":
		return EVENT_STAGE__COMPLETED, nil
	}
}

func ParseEventStageFromLabel(s string) (EventStage, error) {
	switch s {
	default:
		return EVENT_STAGE_UNKNOWN, InvalidEventStage
	case "":
		return EVENT_STAGE_UNKNOWN, nil
	case "RECEIVED":
		return EVENT_STAGE__RECEIVED, nil
	case "HANDLED":
		return EVENT_STAGE__HANDLED, nil
	case "COMPLETED":
		return EVENT_STAGE__COMPLETED, nil
	}
}

func (v EventStage) Int() int {
	return int(v)
}

func (v EventStage) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case EVENT_STAGE_UNKNOWN:
		return ""
	case EVENT_STAGE__RECEIVED:
		return "RECEIVED"
	case EVENT_STAGE__HANDLED:
		return "HANDLED"
	case EVENT_STAGE__COMPLETED:
		return "COMPLETED"
	}
}

func (v EventStage) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case EVENT_STAGE_UNKNOWN:
		return ""
	case EVENT_STAGE__RECEIVED:
		return "RECEIVED"
	case EVENT_STAGE__HANDLED:
		return "HANDLED"
	case EVENT_STAGE__COMPLETED:
		return "COMPLETED"
	}
}

func (v EventStage) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.EventStage"
}

func (v EventStage) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{EVENT_STAGE__RECEIVED, EVENT_STAGE__HANDLED, EVENT_STAGE__COMPLETED}
}

func (v EventStage) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidEventStage
	}
	return []byte(s), nil
}

func (v *EventStage) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseEventStageFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *EventStage) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = EventStage(i)
	return nil
}

func (v EventStage) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
