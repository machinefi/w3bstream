package enums

//go:generate toolkit gen enum Protocol

// Deprecated: Protocol
type Protocol int8

const (
	PROTOCOL_UNKNOWN   Protocol = iota
	PROTOCOL__TCP               // tcp
	PROTOCOL__UDP               // udp
	PROTOCOL__WEBSOCET          // websocket
	PROTOCOL__HTTP              // http
	PROTOCOL__HTTPS             // https
	PROTOCOL__MQTT              // mqtt
)
