package enums

//go:generate toolkit gen enum AccountIdentityType
type AccountIdentityType uint8

const (
	ACCOUNT_IDENTITY_TYPE_UNKNOWN AccountIdentityType = iota
	ACCOUNT_IDENTITY_TYPE__MOBILE
	ACCOUNT_IDENTITY_TYPE__EMAIL
	ACCOUNT_IDENTITY_TYPE__USERNAME
	ACCOUNT_IDENTITY_TYPE__ETHADDRESS
)
