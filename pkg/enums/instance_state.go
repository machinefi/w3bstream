package enums

//go:generate toolkit gen enum InstanceState
type InstanceState uint8

const (
	INSTANCE_STATE_UNKNOWN InstanceState = iota
	INSTANCE_STATE__CREATED
	INSTANCE_STATE__STARTED
	INSTANCE_STATE__STOPPED
	INSTANCE_STATE__LOAD_FAIL
)
