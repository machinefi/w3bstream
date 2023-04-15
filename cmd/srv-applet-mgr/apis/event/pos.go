package event

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/client"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/event"
)

type HandleEvent struct {
	httpx.MethodPost
	Channel              string `in:"path" name:"channel"`
	event.HandleEventReq `in:"body"`
}

func (r *HandleEvent) Path() string { return "/proxy/:channel" }

func (r *HandleEvent) Output(ctx context.Context) (interface{}, error) {
	ev := &eventpb.Event{}
	req, err := http.NewRequest(http.MethodPost, "localhost:8888", bytes.NewBuffer(nil))
	if err != nil {
		return nil, status.EventForwardFailed.StatusErr().WithDesc(err.Error())
	}
	req.Header.Set("Authorization", ev.Header.Token)

	rsp := interface{}(nil)
	if _, err = cli.Do(ctx, req).Into(rsp); err != nil {
		return nil, status.EventForwardFailed.StatusErr().WithDesc(err.Error())
	}
	return rsp, nil
}

var (
	cli = &client.Client{
		Host:    "localhost",
		Port:    8888, // TODO event server
		Timeout: 10 * time.Second,
	}
)

type PublisherTokenAuthProvider struct {
	jwt.Auth
}

type EventHandle struct {
	httpx.MethodPost
	Channel   string `in:"path" name:"channel"`
	EventType string `in:"path" name:"eventType"`
	EventID   string `in:"path" name:"eventID"`
	Payload   []byte `in:"body"`
}

func (r *EventHandle) Path() string {
	return "/transport/:channel/:eventType/:eventID"
}

func (r *EventHandle) Output(ctx context.Context) (interface{}, error) {
	var err error

	p := middleware.MustPublisher(ctx)
	ctx, err = p.WithStrategiesContextByChannelAndType(ctx, r.Channel, r.EventType)
	if err != nil {
		return nil, err
	}

	return event.OnEventReceived(ctx, r.Payload), nil
}
