package mqtt

import (
	me "github.com/iotexproject/w3bstream/pkg/modules/event"
)

type event struct {
	handler     string
	projectID   string
	appletID    string
	publisherID string
	data        []byte
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
