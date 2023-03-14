package account

type EthAccountRegisterReq struct {
	Address EthAddress `json:"address"`
}

//openapi:strfmt eth-address
type EthAddress string

func (v EthAddress) IsZero() bool { return v == "" }

func (v *EthAddress) String() string { return string(*v) }

func (v *EthAddress) UnmarshalText(txt []byte) error {
	*v = EthAddress(txt)
	return nil
}

func (v EthAddress) MarshalText() ([]byte, error) { return []byte(v), nil }
