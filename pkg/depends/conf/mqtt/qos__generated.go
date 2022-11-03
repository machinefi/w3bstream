// This is a generated source file. DO NOT EDIT
// Source: mqtt/qos__generated.go

package mqtt

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidQOS = errors.New("invalid QOS type")

func ParseQOSFromString(s string) (QOS, error) {
	switch s {
	default:
		return QOS_UNKNOWN, InvalidQOS
	case "":
		return QOS_UNKNOWN, nil
	case "ONCE":
		return QOS__ONCE, nil
	case "AT_LEAST_ONCE":
		return QOS__AT_LEAST_ONCE, nil
	case "ONLY_ONCE":
		return QOS__ONLY_ONCE, nil
	}
}

func ParseQOSFromLabel(s string) (QOS, error) {
	switch s {
	default:
		return QOS_UNKNOWN, InvalidQOS
	case "":
		return QOS_UNKNOWN, nil
	case "0":
		return QOS__ONCE, nil
	case "1":
		return QOS__AT_LEAST_ONCE, nil
	case "2":
		return QOS__ONLY_ONCE, nil
	}
}

func (v QOS) Int() int {
	return int(v)
}

func (v QOS) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case QOS_UNKNOWN:
		return ""
	case QOS__ONCE:
		return "ONCE"
	case QOS__AT_LEAST_ONCE:
		return "AT_LEAST_ONCE"
	case QOS__ONLY_ONCE:
		return "ONLY_ONCE"
	}
}

func (v QOS) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case QOS_UNKNOWN:
		return ""
	case QOS__ONCE:
		return "0"
	case QOS__AT_LEAST_ONCE:
		return "1"
	case QOS__ONLY_ONCE:
		return "2"
	}
}

func (v QOS) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt.QOS"
}

func (v QOS) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{QOS__ONCE, QOS__AT_LEAST_ONCE, QOS__ONLY_ONCE}
}

func (v QOS) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidQOS
	}
	return []byte(s), nil
}

func (v *QOS) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseQOSFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *QOS) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = QOS(i)
	return nil
}

func (v QOS) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
