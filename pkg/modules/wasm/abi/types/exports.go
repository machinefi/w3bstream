package types

type Exports interface {
	// OnEventReceived pass data to wasm to handle event
	// entry wasm entry
	// typ event type
	// data payload
	OnEventReceived(entry string, typ string, data []byte) (interface{}, error)
}
