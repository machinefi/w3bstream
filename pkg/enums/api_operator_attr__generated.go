// This is a generated source file. DO NOT EDIT
// Source: enums/api_operator_attr__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidApiOperatorAttr = errors.New("invalid ApiOperatorAttr type")

func ParseApiOperatorAttrFromString(s string) (ApiOperatorAttr, error) {
	switch s {
	default:
		return API_OPERATOR_ATTR_UNKNOWN, InvalidApiOperatorAttr
	case "":
		return API_OPERATOR_ATTR_UNKNOWN, nil
	case "COMMON":
		return API_OPERATOR_ATTR__COMMON, nil
	case "INTERNAL":
		return API_OPERATOR_ATTR__INTERNAL, nil
	case "ADMIN_ONLY":
		return API_OPERATOR_ATTR__ADMIN_ONLY, nil
	case "PUBLIC":
		return API_OPERATOR_ATTR__PUBLIC, nil
	}
}

func ParseApiOperatorAttrFromLabel(s string) (ApiOperatorAttr, error) {
	switch s {
	default:
		return API_OPERATOR_ATTR_UNKNOWN, InvalidApiOperatorAttr
	case "":
		return API_OPERATOR_ATTR_UNKNOWN, nil
	case "common operator for all users with authorization, used for default attr":
		return API_OPERATOR_ATTR__COMMON, nil
	case "internal for debugging or maintaining only":
		return API_OPERATOR_ATTR__INTERNAL, nil
	case "only admin can access":
		return API_OPERATOR_ATTR__ADMIN_ONLY, nil
	case "can access without authorization":
		return API_OPERATOR_ATTR__PUBLIC, nil
	}
}

func (v ApiOperatorAttr) Int() int {
	return int(v)
}

func (v ApiOperatorAttr) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case API_OPERATOR_ATTR_UNKNOWN:
		return ""
	case API_OPERATOR_ATTR__COMMON:
		return "COMMON"
	case API_OPERATOR_ATTR__INTERNAL:
		return "INTERNAL"
	case API_OPERATOR_ATTR__ADMIN_ONLY:
		return "ADMIN_ONLY"
	case API_OPERATOR_ATTR__PUBLIC:
		return "PUBLIC"
	}
}

func (v ApiOperatorAttr) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case API_OPERATOR_ATTR_UNKNOWN:
		return ""
	case API_OPERATOR_ATTR__COMMON:
		return "common operator for all users with authorization, used for default attr"
	case API_OPERATOR_ATTR__INTERNAL:
		return "internal for debugging or maintaining only"
	case API_OPERATOR_ATTR__ADMIN_ONLY:
		return "only admin can access"
	case API_OPERATOR_ATTR__PUBLIC:
		return "can access without authorization"
	}
}

func (v ApiOperatorAttr) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.ApiOperatorAttr"
}

func (v ApiOperatorAttr) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{API_OPERATOR_ATTR__COMMON, API_OPERATOR_ATTR__INTERNAL, API_OPERATOR_ATTR__ADMIN_ONLY, API_OPERATOR_ATTR__PUBLIC}
}

func (v ApiOperatorAttr) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidApiOperatorAttr
	}
	return []byte(s), nil
}

func (v *ApiOperatorAttr) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseApiOperatorAttrFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *ApiOperatorAttr) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = ApiOperatorAttr(i)
	return nil
}

func (v ApiOperatorAttr) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
