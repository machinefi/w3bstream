package enums

//go:generate toolkit gen enum AccountRole
type AccountRole uint8

const (
	ACCOUNT_ROLE_UNKNOWN AccountRole = iota
	ACCOUNT_ROLE__ADMIN
	ACCOUNT_ROLE__DEVELOPER
)
