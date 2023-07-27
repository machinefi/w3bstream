package enums

//go:generate toolkit gen enum FlowOperator
type FlowOperator uint8

const (
	FLOW_OPERATOR_UNKNOWN FlowOperator = iota
	FLOW_OPERATOR__FILTER
	FLOW_OPERATOR__MAP
	FLOW_OPERATOR__FLATMAP
	FLOW_OPERATOR__WINDOW
	FLOW_OPERATOR__GROUP
	FLOW_OPERATOR__REDUCE
)
