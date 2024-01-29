package enums

//go:generate toolkit gen enum EventStage
type EventStage int8

const (
	EVENT_STAGE_UNKNOWN EventStage = iota
	EVENT_STAGE__RECEIVED
	EVENT_STAGE__HANDLED
	EVENT_STAGE__COMPLETED
)
