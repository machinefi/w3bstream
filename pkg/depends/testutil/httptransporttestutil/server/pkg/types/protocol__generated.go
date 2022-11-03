// This is a generated source file. DO NOT EDIT
// Source: types/protocol__generated.go

package types

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidProtocol = errors.New("invalid Protocol type")

func ParseProtocolFromString(s string) (Protocol, error) {
	switch s {
	default:
		return PROTOCOL_UNKNOWN, InvalidProtocol
	case "":
		return PROTOCOL_UNKNOWN, nil
	case "HTTP":
		return PROTOCOL__HTTP, nil
	case "HTTPS":
		return PROTOCOL__HTTPS, nil
	case "TCP":
		return PROTOCOL__TCP, nil
	}
}

func ParseProtocolFromLabel(s string) (Protocol, error) {
	switch s {
	default:
		return PROTOCOL_UNKNOWN, InvalidProtocol
	case "":
		return PROTOCOL_UNKNOWN, nil
	case "http":
		return PROTOCOL__HTTP, nil
	case "https":
		return PROTOCOL__HTTPS, nil
	case "TCP":
		return PROTOCOL__TCP, nil
	}
}

func (v Protocol) Int() int {
	return int(v)
}

func (v Protocol) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case PROTOCOL_UNKNOWN:
		return ""
	case PROTOCOL__HTTP:
		return "HTTP"
	case PROTOCOL__HTTPS:
		return "HTTPS"
	case PROTOCOL__TCP:
		return "TCP"
	}
}

func (v Protocol) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case PROTOCOL_UNKNOWN:
		return ""
	case PROTOCOL__HTTP:
		return "http"
	case PROTOCOL__HTTPS:
		return "https"
	case PROTOCOL__TCP:
		return "TCP"
	}
}

func (v Protocol) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/depends/testutil/httptransporttestutil/server/pkg/types.Protocol"
}

func (v Protocol) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{PROTOCOL__HTTP, PROTOCOL__HTTPS, PROTOCOL__TCP}
}

func (v Protocol) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidProtocol
	}
	return []byte(s), nil
}

func (v *Protocol) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseProtocolFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *Protocol) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = Protocol(i)
	return nil
}

func (v Protocol) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
