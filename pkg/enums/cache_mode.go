package enums

//go:generate toolkit gen enum CacheMode
type CacheMode uint8

const (
	CACHE_MODE_UNKNOWN CacheMode = iota
	CACHE_MODE__MEMORY
	CACHE_MODE__REDIS
)
