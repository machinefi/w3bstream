// This is a generated source file. DO NOT EDIT
// Source: example/scheme__generated.go

package example

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidScheme = errors.New("invalid Scheme type")

func ParseSchemeFromString(s string) (Scheme, error) {
	switch s {
	default:
		return SCHEME_UNKNOWN, InvalidScheme
	case "":
		return SCHEME_UNKNOWN, nil
	case "HTTP":
		return SCHEME__HTTP, nil
	case "HTTPS":
		return SCHEME__HTTPS, nil
	case "TCP":
		return SCHEME__TCP, nil
	case "UDP":
		return SCHEME__UDP, nil
	case "QUIC":
		return SCHEME__QUIC, nil
	}
}

func ParseSchemeFromLabel(s string) (Scheme, error) {
	switch s {
	default:
		return SCHEME_UNKNOWN, InvalidScheme
	case "":
		return SCHEME_UNKNOWN, nil
	case "http":
		return SCHEME__HTTP, nil
	case "https":
		return SCHEME__HTTPS, nil
	case "TCP":
		return SCHEME__TCP, nil
	case "UDP":
		return SCHEME__UDP, nil
	case "QUIC":
		return SCHEME__QUIC, nil
	}
}

func (v Scheme) Int() int {
	return int(v)
}

func (v Scheme) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case SCHEME_UNKNOWN:
		return ""
	case SCHEME__HTTP:
		return "HTTP"
	case SCHEME__HTTPS:
		return "HTTPS"
	case SCHEME__TCP:
		return "TCP"
	case SCHEME__UDP:
		return "UDP"
	case SCHEME__QUIC:
		return "QUIC"
	}
}

func (v Scheme) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case SCHEME_UNKNOWN:
		return ""
	case SCHEME__HTTP:
		return "http"
	case SCHEME__HTTPS:
		return "https"
	case SCHEME__TCP:
		return "TCP"
	case SCHEME__UDP:
		return "UDP"
	case SCHEME__QUIC:
		return "QUIC"
	}
}

func (v Scheme) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/depends/kit/enumgen/__examples__.Scheme"
}

func (v Scheme) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{SCHEME__HTTP, SCHEME__HTTPS, SCHEME__TCP, SCHEME__UDP, SCHEME__QUIC}
}

func (v Scheme) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidScheme
	}
	return []byte(s), nil
}

func (v *Scheme) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseSchemeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *Scheme) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = Scheme(i)
	return nil
}

func (v Scheme) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
