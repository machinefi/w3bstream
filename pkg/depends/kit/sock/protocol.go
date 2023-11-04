package sock

type ProtocolType uint8

const (
	ProtocolUnknown ProtocolType = iota
	ProtocolTCP
	ProtocolUDP
	ProtocolUnix
	ProtocolQUIC // @todo
)

func (p ProtocolType) Network() string {
	switch p {
	case ProtocolTCP:
		return "tcp"
	case ProtocolUDP:
		return "udp"
	case ProtocolUnix:
		return "unix"
	case ProtocolQUIC:
		return "quic"
	default:
		return ""
	}
}
