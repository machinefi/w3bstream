package enums

//go:generate toolkit gen enum FlowSink
type FlowSink uint8

const (
	FLOW_SINK_UNKNOWN FlowSink = iota
	FLOW_SINK__RMDB
	FLOW_SINK__BLOCKCHAIN
)
