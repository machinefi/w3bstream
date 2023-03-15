package enums

//go:generate toolkit gen enum AccountSource
type AccountSource uint8

const (
	ACCOUNT_SOURCE_UNKNOWN AccountSource = iota
	ACCOUNT_SOURCE__INIT
	ACCOUNT_SOURCE__SUBMIT
)
