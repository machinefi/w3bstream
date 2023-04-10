package enums

//go:generate toolkit gen enum ProtoMode

// ProtoMode event protocol adaption mode
type ProtoMode uint8

const (
	PROTO_MODE_UNKNOWN ProtoMode = iota + 0
	PROTO_MODE__PEBBLE
)
