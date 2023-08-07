package enums

//go:generate toolkit gen enum OperatorKeyType
type OperatorKeyType uint8

const (
	OPERATOR_KEY_UNKNOWN OperatorKeyType = iota
	OPERATOR_KEY__ECDSA
	OPERATOR_KEY__ED25519
)
