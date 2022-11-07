// This is a generated source file. DO NOT EDIT
// Source: enums/mqtt_qos__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidMqttQOS = errors.New("invalid MqttQOS type")

func ParseMqttQOSFromString(s string) (MqttQOS, error) {
	switch s {
	default:
		return MQTT_QOS_UNKNOWN, InvalidMqttQOS
	case "":
		return MQTT_QOS_UNKNOWN, nil
	case "ONCE":
		return MQTT_QOS__ONCE, nil
	case "AT_LEAST_ONCE":
		return MQTT_QOS__AT_LEAST_ONCE, nil
	case "ONLY_ONCE":
		return MQTT_QOS__ONLY_ONCE, nil
	}
}

func ParseMqttQOSFromLabel(s string) (MqttQOS, error) {
	switch s {
	default:
		return MQTT_QOS_UNKNOWN, InvalidMqttQOS
	case "":
		return MQTT_QOS_UNKNOWN, nil
	case "0":
		return MQTT_QOS__ONCE, nil
	case "1":
		return MQTT_QOS__AT_LEAST_ONCE, nil
	case "2":
		return MQTT_QOS__ONLY_ONCE, nil
	}
}

func (v MqttQOS) Int() int {
	return int(v)
}

func (v MqttQOS) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case MQTT_QOS_UNKNOWN:
		return ""
	case MQTT_QOS__ONCE:
		return "ONCE"
	case MQTT_QOS__AT_LEAST_ONCE:
		return "AT_LEAST_ONCE"
	case MQTT_QOS__ONLY_ONCE:
		return "ONLY_ONCE"
	}
}

func (v MqttQOS) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case MQTT_QOS_UNKNOWN:
		return ""
	case MQTT_QOS__ONCE:
		return "0"
	case MQTT_QOS__AT_LEAST_ONCE:
		return "1"
	case MQTT_QOS__ONLY_ONCE:
		return "2"
	}
}

func (v MqttQOS) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.MqttQOS"
}

func (v MqttQOS) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{MQTT_QOS__ONCE, MQTT_QOS__AT_LEAST_ONCE, MQTT_QOS__ONLY_ONCE}
}

func (v MqttQOS) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidMqttQOS
	}
	return []byte(s), nil
}

func (v *MqttQOS) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseMqttQOSFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *MqttQOS) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = MqttQOS(i)
	return nil
}

func (v MqttQOS) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
