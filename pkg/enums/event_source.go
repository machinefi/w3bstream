package enums

//go:generate toolkit gen enum EventSource
type EventSource int8

const (
	EVENT_SOURCE_UNKNOWN EventSource = iota
	EVENT_SOURCE__MQTT
	EVENT_SOURCE__HTTP
)
