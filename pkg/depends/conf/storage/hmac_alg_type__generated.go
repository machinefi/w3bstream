// This is a generated source file. DO NOT EDIT
// Source: storage/hmac_alg_type__generated.go

package storage

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidHmacAlgType = errors.New("invalid HmacAlgType type")

func ParseHmacAlgTypeFromString(s string) (HmacAlgType, error) {
	switch s {
	default:
		return HMAC_ALG_TYPE_UNKNOWN, InvalidHmacAlgType
	case "":
		return HMAC_ALG_TYPE_UNKNOWN, nil
	case "MD5":
		return HMAC_ALG_TYPE__MD5, nil
	case "SHA1":
		return HMAC_ALG_TYPE__SHA1, nil
	case "SHA256":
		return HMAC_ALG_TYPE__SHA256, nil
	}
}

func ParseHmacAlgTypeFromLabel(s string) (HmacAlgType, error) {
	switch s {
	default:
		return HMAC_ALG_TYPE_UNKNOWN, InvalidHmacAlgType
	case "":
		return HMAC_ALG_TYPE_UNKNOWN, nil
	case "MD5":
		return HMAC_ALG_TYPE__MD5, nil
	case "SHA1":
		return HMAC_ALG_TYPE__SHA1, nil
	case "SHA256":
		return HMAC_ALG_TYPE__SHA256, nil
	}
}

func (v HmacAlgType) Int() int {
	return int(v)
}

func (v HmacAlgType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case HMAC_ALG_TYPE_UNKNOWN:
		return ""
	case HMAC_ALG_TYPE__MD5:
		return "MD5"
	case HMAC_ALG_TYPE__SHA1:
		return "SHA1"
	case HMAC_ALG_TYPE__SHA256:
		return "SHA256"
	}
}

func (v HmacAlgType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case HMAC_ALG_TYPE_UNKNOWN:
		return ""
	case HMAC_ALG_TYPE__MD5:
		return "MD5"
	case HMAC_ALG_TYPE__SHA1:
		return "SHA1"
	case HMAC_ALG_TYPE__SHA256:
		return "SHA256"
	}
}

func (v HmacAlgType) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/depends/conf/storage.HmacAlgType"
}

func (v HmacAlgType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{HMAC_ALG_TYPE__MD5, HMAC_ALG_TYPE__SHA1, HMAC_ALG_TYPE__SHA256}
}

func (v HmacAlgType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidHmacAlgType
	}
	return []byte(s), nil
}

func (v *HmacAlgType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseHmacAlgTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *HmacAlgType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = HmacAlgType(i)
	return nil
}

func (v HmacAlgType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
