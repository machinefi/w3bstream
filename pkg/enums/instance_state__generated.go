// This is a generated source file. DO NOT EDIT
// Source: enums/instance_state__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidInstanceState = errors.New("invalid InstanceState type")

func ParseInstanceStateFromString(s string) (InstanceState, error) {
	switch s {
	default:
		return INSTANCE_STATE_UNKNOWN, InvalidInstanceState
	case "":
		return INSTANCE_STATE_UNKNOWN, nil
	case "CREATED":
		return INSTANCE_STATE__CREATED, nil
	case "STARTED":
		return INSTANCE_STATE__STARTED, nil
	case "STOPPED":
		return INSTANCE_STATE__STOPPED, nil
	}
}

func ParseInstanceStateFromLabel(s string) (InstanceState, error) {
	switch s {
	default:
		return INSTANCE_STATE_UNKNOWN, InvalidInstanceState
	case "":
		return INSTANCE_STATE_UNKNOWN, nil
	case "CREATED":
		return INSTANCE_STATE__CREATED, nil
	case "STARTED":
		return INSTANCE_STATE__STARTED, nil
	case "STOPPED":
		return INSTANCE_STATE__STOPPED, nil
	}
}

func (v InstanceState) Int() int {
	return int(v)
}

func (v InstanceState) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case INSTANCE_STATE_UNKNOWN:
		return ""
	case INSTANCE_STATE__CREATED:
		return "CREATED"
	case INSTANCE_STATE__STARTED:
		return "STARTED"
	case INSTANCE_STATE__STOPPED:
		return "STOPPED"
	}
}

func (v InstanceState) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case INSTANCE_STATE_UNKNOWN:
		return ""
	case INSTANCE_STATE__CREATED:
		return "CREATED"
	case INSTANCE_STATE__STARTED:
		return "STARTED"
	case INSTANCE_STATE__STOPPED:
		return "STOPPED"
	}
}

func (v InstanceState) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.InstanceState"
}

func (v InstanceState) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{INSTANCE_STATE__CREATED, INSTANCE_STATE__STARTED, INSTANCE_STATE__STOPPED}
}

func (v InstanceState) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidInstanceState
	}
	return []byte(s), nil
}

func (v *InstanceState) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseInstanceStateFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *InstanceState) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = InstanceState(i)
	return nil
}

func (v InstanceState) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
