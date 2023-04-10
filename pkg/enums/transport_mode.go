package enums

//go:generate toolkit gen enum TransportMode

// TransportMode event receiving entry
type TransportMode uint8

const (
	TRANSPORT_MODE_UNKNOWN TransportMode = iota + 0
	TRANSPORT_MODE__MQTT
	TRANSPORT_MODE__HTTP
	TRANSPORT_MODE__WEBSOCKET
)
