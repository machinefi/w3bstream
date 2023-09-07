// This is a generated source file. DO NOT EDIT
// Source: enums/transaction_state__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidTransactionState = errors.New("invalid TransactionState type")

func ParseTransactionStateFromString(s string) (TransactionState, error) {
	switch s {
	default:
		return TRANSACTION_STATE_UNKNOWN, InvalidTransactionState
	case "":
		return TRANSACTION_STATE_UNKNOWN, nil
	case "INIT":
		return TRANSACTION_STATE__INIT, nil
	case "PENDING":
		return TRANSACTION_STATE__PENDING, nil
	case "IN_BLOCK":
		return TRANSACTION_STATE__IN_BLOCK, nil
	case "CONFIRMED":
		return TRANSACTION_STATE__CONFIRMED, nil
	case "FAILED":
		return TRANSACTION_STATE__FAILED, nil
	}
}

func ParseTransactionStateFromLabel(s string) (TransactionState, error) {
	switch s {
	default:
		return TRANSACTION_STATE_UNKNOWN, InvalidTransactionState
	case "":
		return TRANSACTION_STATE_UNKNOWN, nil
	case "INIT":
		return TRANSACTION_STATE__INIT, nil
	case "PENDING":
		return TRANSACTION_STATE__PENDING, nil
	case "IN_BLOCK":
		return TRANSACTION_STATE__IN_BLOCK, nil
	case "CONFIRMED":
		return TRANSACTION_STATE__CONFIRMED, nil
	case "FAILED":
		return TRANSACTION_STATE__FAILED, nil
	}
}

func (v TransactionState) Int() int {
	return int(v)
}

func (v TransactionState) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case TRANSACTION_STATE_UNKNOWN:
		return ""
	case TRANSACTION_STATE__INIT:
		return "INIT"
	case TRANSACTION_STATE__PENDING:
		return "PENDING"
	case TRANSACTION_STATE__IN_BLOCK:
		return "IN_BLOCK"
	case TRANSACTION_STATE__CONFIRMED:
		return "CONFIRMED"
	case TRANSACTION_STATE__FAILED:
		return "FAILED"
	}
}

func (v TransactionState) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case TRANSACTION_STATE_UNKNOWN:
		return ""
	case TRANSACTION_STATE__INIT:
		return "INIT"
	case TRANSACTION_STATE__PENDING:
		return "PENDING"
	case TRANSACTION_STATE__IN_BLOCK:
		return "IN_BLOCK"
	case TRANSACTION_STATE__CONFIRMED:
		return "CONFIRMED"
	case TRANSACTION_STATE__FAILED:
		return "FAILED"
	}
}

func (v TransactionState) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.TransactionState"
}

func (v TransactionState) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{TRANSACTION_STATE__INIT, TRANSACTION_STATE__PENDING, TRANSACTION_STATE__IN_BLOCK, TRANSACTION_STATE__CONFIRMED, TRANSACTION_STATE__FAILED}
}

func (v TransactionState) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidTransactionState
	}
	return []byte(s), nil
}

func (v *TransactionState) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseTransactionStateFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *TransactionState) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = TransactionState(i)
	return nil
}

func (v TransactionState) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
