package event

import (
	"github.com/iotexproject/Bumblebee/base/types"

	me "github.com/iotexproject/w3bstream/pkg/modules/event"
)

type event struct {
	handler     string
	projectID   types.SFID
	appletID    types.SFID
	publisherID types.SFID
	data        []byte
	result      chan me.Result
}

func (e *event) Meta() me.MetaData {
	return me.MetaData{
		PublisherID: e.publisherID,
		Handler:     e.handler,
		AppletID:    e.appletID,
		ProjectID:   e.projectID,
	}
}

func (e *event) Raw() []byte {
	return e.data
}

func (e *event) ResultChan() chan<- me.Result {
	return e.result
}
