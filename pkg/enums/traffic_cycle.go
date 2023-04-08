package enums

//go:generate toolkit gen enum TrafficCycle
type TrafficCycle uint8

const (
	TRAFFIC_CYCLE_UNKNOWN TrafficCycle = iota
	TRAFFIC_CYCLE__MINUTE
	TRAFFIC_CYCLE__HOUR
	TRAFFIC_CYCLE__DAY
	TRAFFIC_CYCLE__WEEK
	TRAFFIC_CYCLE__MONTH
	TRAFFIC_CYCLE__YEAR
)
