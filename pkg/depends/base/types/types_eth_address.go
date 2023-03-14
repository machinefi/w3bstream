package types

import (
	"bytes"
	"strings"
)

//openapi:strfmt eth-address
type EthAddress string

func (v EthAddress) IsZero() bool { return v == "" }

func (v *EthAddress) String() string { return strings.ToLower(string(*v)) }

func (v *EthAddress) UnmarshalText(txt []byte) error {
	*v = EthAddress(bytes.ToLower(txt))
	return nil
}

func (v EthAddress) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}
