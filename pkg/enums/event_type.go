package enums

//go:generate toolkit gen enum EventType
type EventType int

const (
	EVENT_TYPE_UNKNOWN EventType = iota
)

const EVENT_TYPE__ANY = EVENT_TYPE_UNKNOWN + 0x7FFFFFFF // any event type
