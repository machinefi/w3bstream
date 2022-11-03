package types

//go:generate toolkit gen enum Protocol
type Protocol int

const (
	PROTOCOL_UNKNOWN Protocol = iota
	PROTOCOL__HTTP            // http
	PROTOCOL__HTTPS           // https
	_
	_1
	_2
	PROTOCOL__TCP
)

func (Protocol) Offset() int {
	return -4
}
