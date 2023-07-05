package enums

//go:generate toolkit gen enum AccessKeyIdentityType
type AccessKeyIdentityType uint8

const (
	ACCESS_KEY_IDENTITY_TYPE_UNKNOWN AccessKeyIdentityType = iota
	ACCESS_KEY_IDENTITY_TYPE__ACCOUNT
	ACCESS_KEY_IDENTITY_TYPE__PUBLISHER
	_ // ACCESS_KEY_IDENTITY_TYPE__SERVICE
)
