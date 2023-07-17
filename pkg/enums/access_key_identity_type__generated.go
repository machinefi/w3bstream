// This is a generated source file. DO NOT EDIT
// Source: enums/access_key_identity_type__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidAccessKeyIdentityType = errors.New("invalid AccessKeyIdentityType type")

func ParseAccessKeyIdentityTypeFromString(s string) (AccessKeyIdentityType, error) {
	switch s {
	default:
		return ACCESS_KEY_IDENTITY_TYPE_UNKNOWN, InvalidAccessKeyIdentityType
	case "":
		return ACCESS_KEY_IDENTITY_TYPE_UNKNOWN, nil
	case "ACCOUNT":
		return ACCESS_KEY_IDENTITY_TYPE__ACCOUNT, nil
	case "PUBLISHER":
		return ACCESS_KEY_IDENTITY_TYPE__PUBLISHER, nil
	}
}

func ParseAccessKeyIdentityTypeFromLabel(s string) (AccessKeyIdentityType, error) {
	switch s {
	default:
		return ACCESS_KEY_IDENTITY_TYPE_UNKNOWN, InvalidAccessKeyIdentityType
	case "":
		return ACCESS_KEY_IDENTITY_TYPE_UNKNOWN, nil
	case "ACCOUNT":
		return ACCESS_KEY_IDENTITY_TYPE__ACCOUNT, nil
	case "PUBLISHER":
		return ACCESS_KEY_IDENTITY_TYPE__PUBLISHER, nil
	}
}

func (v AccessKeyIdentityType) Int() int {
	return int(v)
}

func (v AccessKeyIdentityType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case ACCESS_KEY_IDENTITY_TYPE_UNKNOWN:
		return ""
	case ACCESS_KEY_IDENTITY_TYPE__ACCOUNT:
		return "ACCOUNT"
	case ACCESS_KEY_IDENTITY_TYPE__PUBLISHER:
		return "PUBLISHER"
	}
}

func (v AccessKeyIdentityType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case ACCESS_KEY_IDENTITY_TYPE_UNKNOWN:
		return ""
	case ACCESS_KEY_IDENTITY_TYPE__ACCOUNT:
		return "ACCOUNT"
	case ACCESS_KEY_IDENTITY_TYPE__PUBLISHER:
		return "PUBLISHER"
	}
}

func (v AccessKeyIdentityType) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.AccessKeyIdentityType"
}

func (v AccessKeyIdentityType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{ACCESS_KEY_IDENTITY_TYPE__ACCOUNT, ACCESS_KEY_IDENTITY_TYPE__PUBLISHER}
}

func (v AccessKeyIdentityType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidAccessKeyIdentityType
	}
	return []byte(s), nil
}

func (v *AccessKeyIdentityType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseAccessKeyIdentityTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *AccessKeyIdentityType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = AccessKeyIdentityType(i)
	return nil
}

func (v AccessKeyIdentityType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
