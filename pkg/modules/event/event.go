// event management

package event

import "github.com/iotexproject/Bumblebee/base/types"

type Event interface {
	Meta() MetaData
	Raw() []byte
}

type MetaData struct {
	PublisherID types.SFID
	ProjectID   types.SFID
	AppletID    types.SFID
	Handler     string // optional
}

type Result struct {
	Success bool
	Data    []byte
}

type EventResult interface {
	ResultChan() chan<- Result
}
