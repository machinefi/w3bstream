package example

type Scheme int

const (
	SCHEME_UNKNOWN Scheme = iota
	SCHEME__HTTP          // http
	SCHEME__HTTPS         // https
	_
	_1
	_2
	SCHEME__TCP
	SCHEME__UDP
	SCHEME__QUIC
)

func (Scheme) Offset() int { return 1 }
