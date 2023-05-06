package enums

// InstanceState presents if wasm was loaded to memory and if it can receive data
//
//go:generate toolkit gen enum InstanceState
type InstanceState uint8

const (
	INSTANCE_STATE_UNKNOWN InstanceState = iota
	_
	INSTANCE_STATE__STARTED // ready to receive data
	INSTANCE_STATE__STOPPED // stopped to receive data
)
