// This is a generated source file. DO NOT EDIT
// Source: enums/file_system_mode__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidFileSystemMode = errors.New("invalid FileSystemMode type")

func ParseFileSystemModeFromString(s string) (FileSystemMode, error) {
	switch s {
	default:
		return FILE_SYSTEM_MODE_UNKNOWN, InvalidFileSystemMode
	case "":
		return FILE_SYSTEM_MODE_UNKNOWN, nil
	case "LOCAL":
		return FILE_SYSTEM_MODE__LOCAL, nil
	case "S3":
		return FILE_SYSTEM_MODE__S3, nil
	}
}

func ParseFileSystemModeFromLabel(s string) (FileSystemMode, error) {
	switch s {
	default:
		return FILE_SYSTEM_MODE_UNKNOWN, InvalidFileSystemMode
	case "":
		return FILE_SYSTEM_MODE_UNKNOWN, nil
	case "LOCAL":
		return FILE_SYSTEM_MODE__LOCAL, nil
	case "S3":
		return FILE_SYSTEM_MODE__S3, nil
	}
}

func (v FileSystemMode) Int() int {
	return int(v)
}

func (v FileSystemMode) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case FILE_SYSTEM_MODE_UNKNOWN:
		return ""
	case FILE_SYSTEM_MODE__LOCAL:
		return "LOCAL"
	case FILE_SYSTEM_MODE__S3:
		return "S3"
	}
}

func (v FileSystemMode) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case FILE_SYSTEM_MODE_UNKNOWN:
		return ""
	case FILE_SYSTEM_MODE__LOCAL:
		return "LOCAL"
	case FILE_SYSTEM_MODE__S3:
		return "S3"
	}
}

func (v FileSystemMode) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.FileSystemMode"
}

func (v FileSystemMode) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{FILE_SYSTEM_MODE__LOCAL, FILE_SYSTEM_MODE__S3}
}

func (v FileSystemMode) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidFileSystemMode
	}
	return []byte(s), nil
}

func (v *FileSystemMode) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseFileSystemModeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *FileSystemMode) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = FileSystemMode(i)
	return nil
}

func (v FileSystemMode) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
