// This is a generated source file. DO NOT EDIT
// Source: storage/storage_type__generated.go

package storage

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidStorageType = errors.New("invalid StorageType type")

func ParseStorageTypeFromString(s string) (StorageType, error) {
	switch s {
	default:
		return STORAGE_TYPE_UNKNOWN, InvalidStorageType
	case "":
		return STORAGE_TYPE_UNKNOWN, nil
	case "S3":
		return STORAGE_TYPE__S3, nil
	case "FILESYSTEM":
		return STORAGE_TYPE__FILESYSTEM, nil
	case "IPFS":
		return STORAGE_TYPE__IPFS, nil
	}
}

func ParseStorageTypeFromLabel(s string) (StorageType, error) {
	switch s {
	default:
		return STORAGE_TYPE_UNKNOWN, InvalidStorageType
	case "":
		return STORAGE_TYPE_UNKNOWN, nil
	case "S3":
		return STORAGE_TYPE__S3, nil
	case "FILESYSTEM":
		return STORAGE_TYPE__FILESYSTEM, nil
	case "IPFS":
		return STORAGE_TYPE__IPFS, nil
	}
}

func (v StorageType) Int() int {
	return int(v)
}

func (v StorageType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case STORAGE_TYPE_UNKNOWN:
		return ""
	case STORAGE_TYPE__S3:
		return "S3"
	case STORAGE_TYPE__FILESYSTEM:
		return "FILESYSTEM"
	case STORAGE_TYPE__IPFS:
		return "IPFS"
	}
}

func (v StorageType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case STORAGE_TYPE_UNKNOWN:
		return ""
	case STORAGE_TYPE__S3:
		return "S3"
	case STORAGE_TYPE__FILESYSTEM:
		return "FILESYSTEM"
	case STORAGE_TYPE__IPFS:
		return "IPFS"
	}
}

func (v StorageType) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/depends/conf/storage.StorageType"
}

func (v StorageType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{STORAGE_TYPE__S3, STORAGE_TYPE__FILESYSTEM, STORAGE_TYPE__IPFS}
}

func (v StorageType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidStorageType
	}
	return []byte(s), nil
}

func (v *StorageType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseStorageTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *StorageType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = StorageType(i)
	return nil
}

func (v StorageType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
