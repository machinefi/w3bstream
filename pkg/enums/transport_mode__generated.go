// This is a generated source file. DO NOT EDIT
// Source: enums/transport_mode__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidTransportMode = errors.New("invalid TransportMode type")

func ParseTransportModeFromString(s string) (TransportMode, error) {
	switch s {
	default:
		return TRANSPORT_MODE_UNKNOWN, InvalidTransportMode
	case "":
		return TRANSPORT_MODE_UNKNOWN, nil
	case "MQTT":
		return TRANSPORT_MODE__MQTT, nil
	case "HTTP":
		return TRANSPORT_MODE__HTTP, nil
	case "WEBSOCKET":
		return TRANSPORT_MODE__WEBSOCKET, nil
	}
}

func ParseTransportModeFromLabel(s string) (TransportMode, error) {
	switch s {
	default:
		return TRANSPORT_MODE_UNKNOWN, InvalidTransportMode
	case "":
		return TRANSPORT_MODE_UNKNOWN, nil
	case "MQTT":
		return TRANSPORT_MODE__MQTT, nil
	case "HTTP":
		return TRANSPORT_MODE__HTTP, nil
	case "WEBSOCKET":
		return TRANSPORT_MODE__WEBSOCKET, nil
	}
}

func (v TransportMode) Int() int {
	return int(v)
}

func (v TransportMode) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case TRANSPORT_MODE_UNKNOWN:
		return ""
	case TRANSPORT_MODE__MQTT:
		return "MQTT"
	case TRANSPORT_MODE__HTTP:
		return "HTTP"
	case TRANSPORT_MODE__WEBSOCKET:
		return "WEBSOCKET"
	}
}

func (v TransportMode) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case TRANSPORT_MODE_UNKNOWN:
		return ""
	case TRANSPORT_MODE__MQTT:
		return "MQTT"
	case TRANSPORT_MODE__HTTP:
		return "HTTP"
	case TRANSPORT_MODE__WEBSOCKET:
		return "WEBSOCKET"
	}
}

func (v TransportMode) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.TransportMode"
}

func (v TransportMode) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{TRANSPORT_MODE__MQTT, TRANSPORT_MODE__HTTP, TRANSPORT_MODE__WEBSOCKET}
}

func (v TransportMode) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidTransportMode
	}
	return []byte(s), nil
}

func (v *TransportMode) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseTransportModeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *TransportMode) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = TransportMode(i)
	return nil
}

func (v TransportMode) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
