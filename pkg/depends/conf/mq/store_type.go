package mq

//go:generate toolkit gen enum StoreType
type StoreType uint8

const (
	STORE_TYPE_UNKNOWN StoreType = iota
	STORE_TYPE__MEM
	STORE_TYPE__REDIS
)
