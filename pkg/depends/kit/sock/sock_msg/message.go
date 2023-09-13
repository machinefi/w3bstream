package sock_msg

import (
	"github.com/pkg/errors"

	qbuf "github.com/machinefi/w3bstream/pkg/depends/kit/sock/sock_buf"
)

type Type interface {
	String() string
}

type ID interface {
	String() string
}

type Message interface {
	ID() ID
	Type() Type
}

type NamedMessage interface {
	Name() string
}

type WithTimestamp interface {
	GetTimestamp() int64
}

type WithErrorCheck interface {
	Err() error
}

type Parser interface {
	Marshal(qbuf.Buffer, Message) error
	Unmarshal(qbuf.Buffer) (Message, error)
}

var ErrParseTCPDataLack = errors.New("data lack")
