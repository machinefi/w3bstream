package sock

import "fmt"

type Errno int

const (
	ENodeTimeout Errno = 0x01 + iota
	ENodeContextDone
	ENodeMessage // received nil message,
	ENodeMarshal
	ENodeUnmarshal
	ENodeClosed
	ENodeDial
	ENodeRead
	ENodeWrite
	ENodeListen
	ENodeResolve
	ENodeInvalidParser
	ENodeInvalidRemoteAddr
	ENodeInvalidProtocol
	ENodeInvalidListenAddr
	ENodeOption
	EMessageUnbound
	EMessageTimeout
	EMessageIdRepeated
)

var es = [...]string{
	ENodeTimeout:           "sock.Timeout",
	ENodeContextDone:       "sock.Context",
	ENodeMessage:           "sock.Message",
	ENodeMarshal:           "sock.Marshal",
	ENodeUnmarshal:         "sock.Unmarshal",
	ENodeClosed:            "sock.Closed",
	ENodeDial:              "sock.Dial",
	ENodeRead:              "sock.Read",
	ENodeWrite:             "sock.Write",
	ENodeListen:            "sock.Listen",
	ENodeResolve:           "sock.Resolve",
	ENodeInvalidParser:     "sock.InvalidParser",
	ENodeInvalidRemoteAddr: "sock.InvalidRemote",
	ENodeInvalidProtocol:   "sock.InvalidProtocol",
	ENodeInvalidListenAddr: "sock.InvalidListen",
	ENodeOption:            "sock.Option",
	EMessageUnbound:        "sock.UnboundMessage",
	EMessageTimeout:        "sock.MessageTimeout",
	EMessageIdRepeated:     "sock.MessageIdRepeated",
}

func (e Errno) Error() string {
	if e >= ENodeTimeout && e <= EMessageIdRepeated {
		return es[e]
	}
	return ""
}

func (e Errno) WithMessage(msg string) *ErrnoWithMsg {
	return &ErrnoWithMsg{
		Errno: e,
		msg:   msg,
	}
}

func (e Errno) WithError(err error) *ErrnoWithErr {
	return &ErrnoWithErr{
		Errno: e,
		err:   err,
	}
}

type ErrnoWithMsg struct {
	Errno
	msg string
}

func (e *ErrnoWithMsg) Error() string {
	return fmt.Sprintf("%v: %s", e.Errno, e.msg)
}

func (e *ErrnoWithMsg) Unwrap() error {
	return e.Errno
}

type ErrnoWithErr struct {
	Errno
	err error
}

func (e *ErrnoWithErr) Error() string {
	return fmt.Sprintf("%v: %s", e.Errno, e.err.Error())
}

func (e *ErrnoWithErr) Unwrap() error {
	return e.Errno
}

func IsTimeoutError(e error) bool {
	switch v := e.(type) {
	case Errno:
		return v == EMessageTimeout || v == ENodeTimeout
	case *ErrnoWithMsg:
		return v.Errno == EMessageTimeout || v.Errno == ENodeTimeout
	case *ErrnoWithErr:
		return v.Errno == EMessageTimeout || v.Errno == ENodeTimeout
	default:
		return false
	}
}

func IsNodeClosedError(e error) bool {
	switch v := e.(type) {
	case Errno:
		return v == ENodeClosed
	case *ErrnoWithMsg:
		return v.Errno == ENodeClosed
	case *ErrnoWithErr:
		return v.Errno == ENodeClosed
	default:
		return false
	}
}
