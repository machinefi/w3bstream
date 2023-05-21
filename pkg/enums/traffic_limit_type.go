package enums

//go:generate toolkit gen enum TrafficLimitType
type TrafficLimitType uint8

const (
	TRAFFIC_LIMIT_TYPE_UNKNOWN TrafficLimitType = iota
	TRAFFIC_LIMIT_TYPE__EVENT
	TRAFFIC_LIMIT_TYPE__BLOCKCHAIN
)
