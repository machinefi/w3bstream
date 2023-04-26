package event

import (
	"time"

	"github.com/google/uuid"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/types"
)

type EventReq struct {
	// Channel message channel named (intact project name)
	Channel string `in:"path"  name:"channel"`
	// EventType used for filter strategies created in w3b before
	EventType string `in:"query" name:"eventType,omitempty"`
	// EventID unique id for tracing event under channel
	EventID string `in:"query" name:"eventID,omitempty"`
	// Timestamp event time when publisher do send
	Timestamp int64 `in:"query" name:"timestamp,omitempty"`
	// Payload event payload (binary only)
	Payload []byte `in:"body"`
}

func (r *EventReq) SetDefault() {
	if r.EventType == "" {
		r.EventType = enums.EVENTTYPEDEFAULT
	}
	if r.EventID == "" {
		r.EventID = uuid.NewString() + "_w3b" // flag generated by w3b node
	}
	if r.Timestamp == 0 {
		r.Timestamp = time.Now().UTC().Unix()
	}
}

type Result struct {
	// AppletName applet name(unique) under published channel(project)
	AppletName string `json:"appletName"`
	// InstanceID the unique wasm vm  id
	InstanceID types.SFID `json:"instanceID"`
	// Handler invoked wasm entry name
	Handler string `json:"handler"`
	// ReturnValue wasm call returned value
	ReturnValue []byte `json:"returnValue"`
	// ReturnCode wasm call returned code
	ReturnCode int `json:"code"`
	// Error message instance module, presents result for wasm invoking
	Error string `json:"error,omitempty"`
}

type EventRsp struct {
	// Channel intact project name
	Channel string `json:"channel"`
	// PublisherID publisher(device) unique id in w3b node
	PublisherID types.SFID `json:"publisherID"`
	// EventID same as EventReq.EventID
	EventID string `json:"eventID"`
	// Results result for each wasm invoke, which hits strategies.
	Results []*Result `json:"results"`
	// Error error message from w3b node (api level), different from Result.Error
	Error string `json:"error,omitempty"`
}
