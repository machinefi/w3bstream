package enums

//go:generate toolkit gen enum RateLimitApiType
type RateLimitApiType uint8

const (
	RATE_LIMIT_API_TYPE_UNKNOWN RateLimitApiType = iota
	RATE_LIMIT_API_TYPE__EVENT
	RATE_LIMIT_API_TYPE__BLOCKCHAIN
)
