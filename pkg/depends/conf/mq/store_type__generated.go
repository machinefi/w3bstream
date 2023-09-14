// This is a generated source file. DO NOT EDIT
// Source: mq/store_type__generated.go

package mq

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidStoreType = errors.New("invalid StoreType type")

func ParseStoreTypeFromString(s string) (StoreType, error) {
	switch s {
	default:
		return STORE_TYPE_UNKNOWN, InvalidStoreType
	case "":
		return STORE_TYPE_UNKNOWN, nil
	case "MEM":
		return STORE_TYPE__MEM, nil
	case "REDIS":
		return STORE_TYPE__REDIS, nil
	}
}

func ParseStoreTypeFromLabel(s string) (StoreType, error) {
	switch s {
	default:
		return STORE_TYPE_UNKNOWN, InvalidStoreType
	case "":
		return STORE_TYPE_UNKNOWN, nil
	case "MEM":
		return STORE_TYPE__MEM, nil
	case "REDIS":
		return STORE_TYPE__REDIS, nil
	}
}

func (v StoreType) Int() int {
	return int(v)
}

func (v StoreType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case STORE_TYPE_UNKNOWN:
		return ""
	case STORE_TYPE__MEM:
		return "MEM"
	case STORE_TYPE__REDIS:
		return "REDIS"
	}
}

func (v StoreType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case STORE_TYPE_UNKNOWN:
		return ""
	case STORE_TYPE__MEM:
		return "MEM"
	case STORE_TYPE__REDIS:
		return "REDIS"
	}
}

func (v StoreType) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/depends/conf/mq.StoreType"
}

func (v StoreType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{STORE_TYPE__MEM, STORE_TYPE__REDIS}
}

func (v StoreType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidStoreType
	}
	return []byte(s), nil
}

func (v *StoreType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseStoreTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *StoreType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = StoreType(i)
	return nil
}

func (v StoreType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
