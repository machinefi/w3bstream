// This is a generated source file. DO NOT EDIT
// Source: enums/deploy_cmd__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidDeployCmd = errors.New("invalid DeployCmd type")

func ParseDeployCmdFromString(s string) (DeployCmd, error) {
	switch s {
	default:
		return DEPLOY_CMD_UNKNOWN, InvalidDeployCmd
	case "":
		return DEPLOY_CMD_UNKNOWN, nil
	case "START":
		return DEPLOY_CMD__START, nil
	case "HUNGUP":
		return DEPLOY_CMD__HUNGUP, nil
	}
}

func ParseDeployCmdFromLabel(s string) (DeployCmd, error) {
	switch s {
	default:
		return DEPLOY_CMD_UNKNOWN, InvalidDeployCmd
	case "":
		return DEPLOY_CMD_UNKNOWN, nil
	case "start wasm vm":
		return DEPLOY_CMD__START, nil
	case "stop wasm vm":
		return DEPLOY_CMD__HUNGUP, nil
	}
}

func (v DeployCmd) Int() int {
	return int(v)
}

func (v DeployCmd) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case DEPLOY_CMD_UNKNOWN:
		return ""
	case DEPLOY_CMD__START:
		return "START"
	case DEPLOY_CMD__HUNGUP:
		return "HUNGUP"
	}
}

func (v DeployCmd) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case DEPLOY_CMD_UNKNOWN:
		return ""
	case DEPLOY_CMD__START:
		return "start wasm vm"
	case DEPLOY_CMD__HUNGUP:
		return "stop wasm vm"
	}
}

func (v DeployCmd) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.DeployCmd"
}

func (v DeployCmd) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{DEPLOY_CMD__START, DEPLOY_CMD__HUNGUP}
}

func (v DeployCmd) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidDeployCmd
	}
	return []byte(s), nil
}

func (v *DeployCmd) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseDeployCmdFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *DeployCmd) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = DeployCmd(i)
	return nil
}

func (v DeployCmd) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
